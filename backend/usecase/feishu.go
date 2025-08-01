package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/88250/lute"
	"github.com/Wsine/feishu2md/core"
	fershu2mdUtil "github.com/Wsine/feishu2md/utils"
	"github.com/google/uuid"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdrive1 "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
	larksheets "github.com/larksuite/oapi-sdk-go/v3/service/sheets/v3"
	larkwiki1 "github.com/larksuite/oapi-sdk-go/v3/service/wiki/v1"
	larkwiki2 "github.com/larksuite/oapi-sdk-go/v3/service/wiki/v2"
	"github.com/minio/minio-go/v7"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"

	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/store/s3"
	"github.com/chaitin/panda-wiki/utils"
)

type FeishuUseCase struct {
	logger      *log.Logger
	minioClient *s3.MinioClient
	crawler     *CrawlerUsecase
}

func NewFeishuUseCase(logger *log.Logger, minio *s3.MinioClient, crawler *CrawlerUsecase) *FeishuUseCase {
	return &FeishuUseCase{
		logger:      logger.WithModule("usecase.feishuusecase"),
		minioClient: minio,
		crawler:     crawler,
	}
}

func (f *FeishuUseCase) GetSpacelist(ctx context.Context, req *domain.GetSpaceListReq) ([]*domain.GetSpaceListResp, error) {
	client := lark.NewClient(req.AppID, req.AppSecret)
	var (
		respData  []*domain.GetSpaceListResp
		pageToken string
	)
	for {
		// 创建请求对象
		r := larkwiki2.NewListSpaceReqBuilder().
			PageSize(20).
			PageToken(pageToken).
			Lang(`zh`).
			Build()

		f.logger.Debug("token", log.String("token", req.UserAccessToken))
		// 发起请求
		resp, err := client.Wiki.V2.Space.List(
			ctx,
			r,
			larkcore.WithUserAccessToken(req.UserAccessToken))
		if err != nil {
			f.logger.Error("list space failed", log.Error(err))
			return nil, fmt.Errorf("list space failed: %v", err)
		}

		if resp.Msg != "success" {
			return nil, fmt.Errorf("list space failed: %s", resp.Msg)
		}
		f.logger.Debug("resp", log.Any("resp", resp))
		for _, v := range resp.Data.Items {
			respData = append(respData, &domain.GetSpaceListResp{
				Name:    *v.Name,
				SpaceId: *v.SpaceId,
			})
		}
		if !*resp.Data.HasMore {
			break
		}
		pageToken = *resp.Data.PageToken
	}
	return respData, nil
}

func (f *FeishuUseCase) SearchWiki(ctx context.Context, req *domain.SearchWikiReq) ([]*domain.SearchWikiResp, error) {
	client := lark.NewClient(req.AppID, req.AppSecret)
	var (
		pageToken string
		respData  []*domain.SearchWikiResp
	)
	for {
		// 创建请求对象
		r := larkwiki1.NewSearchNodeReqBuilder().
			PageSize(20).
			PageToken(pageToken).
			Body(larkwiki1.NewSearchNodeReqBodyBuilder().
				Query(req.Query).
				SpaceId(req.SpaceId).
				Build()).
			Build()

		f.logger.Debug("token", log.String("token", req.UserAccessToken))
		// 发起请求
		resp, err := client.Wiki.V1.Node.Search(ctx, r, larkcore.WithUserAccessToken(req.UserAccessToken))
		if err != nil {
			f.logger.Error("search Wiki failed", log.Error(err))
			return nil, fmt.Errorf("search Wiki failed: %v", err)
		}

		if resp.Msg != "success" {
			return nil, fmt.Errorf("search Wiki failed: %s", resp.Msg)
		}
		for _, v := range resp.Data.Items {
			if *v.ObjType == 9 {
				continue
			}
			respData = append(respData, &domain.SearchWikiResp{
				Title:    *v.Title,
				Url:      *v.Url,
				SpaceId:  *v.SpaceId,
				ObjToken: *v.ObjToken,
				ObjType:  *v.ObjType,
			})
			if v.Title == nil || *v.Title == "" {
				f.logger.Info("title is empty", log.Int("obj_type", *v.ObjType))
			}
			f.logger.Info("found wiki", log.String("title", *v.Title), log.Int("obj_type", *v.ObjType))
		}
		if !*resp.Data.HasMore {
			break
		}
		pageToken = *resp.Data.PageToken
	}
	return respData, nil
}

func (f *FeishuUseCase) GetDoc(ctx context.Context, req *domain.GetDocxReq) ([]*domain.GetDocxResp, error) {
	results := parallel.Map(req.Sources, func(source domain.Source, _ int) *domain.GetDocxResp {
		var (
			title   string
			content string
			err     error
		)
		switch source.ObjType {
		case 8: // docx
			f.logger.Info("download [type 8] docx ", log.String("url", source.Url), log.String("token", source.ObjToken))
			title, content, err = f.downloadDocument(ctx, req.AppID, req.AppSecret, source.Url, req.KBID)
			if err != nil {
				f.logger.Error("download [type 8] docx failed", log.Error(err))
				return nil // 返回 nil 表示失败
			}
		case 5: // file
			f.logger.Info("download [type 5]file ", log.String("url", source.Url), log.String("token", source.ObjToken))
			title, content, err = f.downloadFile(ctx, req.AppID, req.AppSecret, source.ObjToken, req.KBID)
			if err != nil {
				f.logger.Error("download [type 5]file failed", log.Error(err))
			}
		case 2: // sheet
			f.logger.Info("download [type 2] sheet", log.String("url", source.Url), log.String("token", source.ObjToken))
			title, content, err = f.downloadSheet(ctx, req.AppID, req.AppSecret, source.ObjToken, req.UserAccessToken, req.KBID)
			if err != nil {
				f.logger.Error("download [type 2] sheet failed", log.Error(err))
			}
		default: // 其他类型
			f.logger.Error("unsupported obj type", log.Int("type", source.ObjType))
			return nil // 返回 nil 表示失败
		}

		return &domain.GetDocxResp{
			Title:   title,
			Content: content,
		}
	})

	// 过滤掉 nil（失败的情况）
	successResults := lo.Filter(results, func(resp *domain.GetDocxResp, _ int) bool {
		return resp != nil
	})

	// 如果有失败的请求，返回错误
	if len(successResults) < len(req.Sources) {
		return successResults, fmt.Errorf("some downloads failed (%d/%d succeeded)", len(successResults), len(req.Sources))
	}
	return successResults, nil
}

func (f *FeishuUseCase) ListDocx(ctx context.Context, req *domain.SearchDocxReq) ([]*domain.SearchDocxResp, error) {
	client := lark.NewClient(req.AppID, req.AppSecret)

	// 使用context.WithCancel实现优雅退出
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error, 1)                     // 缓冲通道防止goroutine泄漏
	dataChan := make(chan *domain.SearchDocxResp, 100) // 适当缓冲提高性能

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		searchDocx(ctx, client, req.UserAccessToken, "", dataChan, errChan)
	}()

	// 单独goroutine等待所有搜索完成
	go func() {
		wg.Wait()
		close(dataChan)
		close(errChan)
	}()

	results := make([]*domain.SearchDocxResp, 0)

	// 处理结果和错误
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case data, ok := <-dataChan:
			if !ok {
				return results, nil
			}
			results = append(results, data)
		case err, ok := <-errChan:
			if ok {
				return nil, err
			}
		}
	}
}

func searchDocx(ctx context.Context, client *lark.Client, accessToken, folderToken string,
	dataChan chan<- *domain.SearchDocxResp, errChan chan<- error,
) {
	// 构建请求
	builder := larkdrive1.NewListFileReqBuilder().
		OrderBy("EditedTime").
		Direction("DESC")

	if folderToken != "" {
		builder.FolderToken(folderToken)
	}

	// 发送请求
	resp, err := client.Drive.V1.File.List(
		ctx, // 使用传入的context
		builder.Build(),
		larkcore.WithUserAccessToken(accessToken),
	)
	if resp.Msg != "success" {
		select {
		case errChan <- fmt.Errorf("search doc failed: %s", resp.Msg):
		case <-ctx.Done():
		}
	}
	if err != nil {

		select {
		case errChan <- fmt.Errorf("search doc failed: %w", err):
		case <-ctx.Done():
		}
		return
	}

	var wg sync.WaitGroup

	for _, v := range resp.Data.Files {
		select {
		case <-ctx.Done():
			return
		default:
			if *v.Type == "folder" {
				wg.Add(1)
				go func(token string) {
					defer wg.Done()
					searchDocx(ctx, client, accessToken, token, dataChan, errChan)
				}(*v.Token)
			} else {
				var temp int
				switch *v.Type {
				case "doc", "docx":
					temp = 8
				case "sheet":
					temp = 2
				case "file":
					temp = 5
				default:
					temp = 0
				}
				select {
				case dataChan <- &domain.SearchDocxResp{
					Name:     *v.Name,
					Url:      *v.Url,
					ObjToken: *v.Token,
					ObjType:  temp,
				}:
				case <-ctx.Done():
				}
			}
		}
	}

	wg.Wait()
}

func (f *FeishuUseCase) downloadDocument(ctx context.Context, appID, secret, url, kbID string) (string, string, error) {
	client := core.NewClient(
		appID, secret,
	)
	// Validate the url to download
	var dlConfig core.Config
	dlConfig.Output = core.OutputConfig{
		ImageDir:        "./temp/images",
		SkipImgDownload: false,
		TitleAsFilename: false,
		UseHTMLTags:     true,
	}
	// 清理图像临时存储目录
	defer func() {
		os.RemoveAll(dlConfig.Output.ImageDir)
	}()
	docType, docToken, err := fershu2mdUtil.ValidateDocumentURL(url)
	if err != nil {
		return "", "", err
	}
	if docType == "wiki" {
		node, err := client.GetWikiNodeInfo(ctx, docToken)
		if err != nil {
			return "", "", fmt.Errorf("GetWikiNodeInfo err: %v for %v", err, url)
		}
		docToken = node.ObjToken
	}

	// Process the download
	docx, blocks, err := client.GetDocxContent(ctx, docToken)
	if err != nil {
		return "", "", err
	}

	parser := core.NewParser(dlConfig.Output)

	title := docx.Title
	if title == "" {
		title = "未命名文档"
	}
	markdown := parser.ParseDocxContent(docx, blocks)
	type imgReplace struct {
		token string
		path  string
	}

	replaceChan := make(chan imgReplace, len(parser.ImgTokens))
	errChan := make(chan error, len(parser.ImgTokens))
	var wg sync.WaitGroup
	for _, imgToken := range parser.ImgTokens {
		wg.Add(1)
		go func(token string) {
			defer wg.Done()

			localLink, err := client.DownloadImage(ctx, token, dlConfig.Output.ImageDir)
			if err != nil {
				errChan <- fmt.Errorf("下载图片失败: %v", err)
				return
			}

			osspath, err := utils.UploadImage(ctx, f.minioClient, localLink, kbID)
			if err != nil {
				errChan <- fmt.Errorf("上传图片失败: %v", err)
				return
			}

			replaceChan <- imgReplace{token: token, path: osspath}
		}(imgToken)
	}
	go func() {
		wg.Wait()
		close(replaceChan)
		close(errChan)
	}()

	// 检查错误
	select {
	case err := <-errChan:
		return "", "", err
	default:
	}
	replaces := make(map[string]string)
	for r := range replaceChan {
		replaces[r.token] = r.path
	}
	replacePairs := make([]string, 0, len(replaces)*2)
	for token, path := range replaces {
		replacePairs = append(replacePairs, token, path)
	}

	// 创建一次性替换器
	replacer := strings.NewReplacer(replacePairs...)
	markdown = replacer.Replace(markdown)

	// Format the markdown document
	engine := lute.New(func(l *lute.Lute) {
		l.RenderOptions.AutoSpace = true
	})
	result := engine.FormatStr("md", markdown)
	return title, result, nil
}

func (f *FeishuUseCase) downloadFile(ctx context.Context, appID, secret, fileToken, kbID string) (string, string, error) {
	client := lark.NewClient(
		appID, secret,
	)

	req := larkdrive1.NewDownloadFileReqBuilder().
		FileToken(fileToken).
		Build()

	// 发起请求
	resp, err := client.Drive.V1.File.Download(ctx, req)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return "", "", fmt.Errorf("下载文件失败: %s", resp.CodeError.Msg)
	}

	// 业务处理
	var contentType string
	ext := strings.ToLower(filepath.Ext(resp.FileName))

	switch ext {
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	case ".pdf":
		contentType = "application/pdf"
	case ".doc", ".docx":
		contentType = "application/msword"
	case ".xls", ".xlsx":
		contentType = "application/vnd.ms-excel"
	case ".ppt", ".pptx":
		contentType = "application/vnd.ms-powerpoint"
	case ".txt":
		contentType = "text/plain"
	case ".html", ".htm":
		contentType = "text/html"
	default:
		contentType = "application/octet-stream" // 未知类型
	}

	imgName := fmt.Sprintf("%s/%s%s", kbID, uuid.New().String(), ext)

	buf, _ := io.ReadAll(resp.File)

	_, err = f.minioClient.PutObject(
		ctx,
		domain.Bucket,
		imgName,
		bytes.NewReader(buf),
		int64(len(buf)),
		minio.PutObjectOptions{
			ContentType: contentType,
			UserMetadata: map[string]string{
				"originalname": resp.FileName,
			},
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("上传文件到OSS失败: %s", err)
	}
	ossPath := fmt.Sprintf("/%s/%s", domain.Bucket, imgName)
	res, err := f.crawler.ScrapeURL(ctx, ossPath, kbID)
	if err != nil {
		return "", "", err
	}
	name := resp.FileName
	if name == "" {
		name = res.Title
	}
	if name == "" {
		name = "未命名文件"
	}
	return name, res.Content, nil
}

func getSheetInfo(ctx context.Context, appID, secret, sheetToken, usserAccessToken string) (*larksheets.GetSpreadsheet, error) {
	client := lark.NewClient(appID, secret)
	// 创建请求对象
	req := larksheets.NewGetSpreadsheetReqBuilder().
		SpreadsheetToken(sheetToken).
		Build()

	// 发起请求
	resp, err := client.Sheets.V3.Spreadsheet.Get(ctx, req, larkcore.WithUserAccessToken(usserAccessToken))
	if err != nil {
		return nil, err
	}
	return resp.Data.Spreadsheet, nil
}

func (f *FeishuUseCase) downloadSheet(ctx context.Context, appID, secret, sheetToken, usserAccessToken, kbID string) (string, string, error) {
	sheets, err := f.getSheets(ctx, appID, secret, sheetToken, usserAccessToken)
	if err != nil {
		f.logger.Error("get sheets failed", log.Error(err))
		return "", "", err
	}
	for i, sheet := range sheets {
		f.logger.Debug("Sheet details",
			log.Int("index", i),
			log.String("sheetId", *sheet.SheetId),
			log.String("title", *sheet.Title),
			log.Int("index", *sheet.Index)) // 如果有这个字段
	}
	subIDs := []string{*sheets[0].SheetId}
	tickits, err := f.creatExportTask(ctx, appID, secret, sheetToken, usserAccessToken, subIDs)
	if err != nil {
		f.logger.Error("export task failed", log.Error(err))
		return "", "", err
	}
	for i, tickit := range tickits {
		f.logger.Debug("Export task details",
			log.Int("index", i),
			log.String("ticket", tickit))
	}

	getExportTaskResps := parallel.Map(tickits, func(tickit string, _ int) *getExportTaskResp {
		ticker := time.NewTicker(10 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				res, err := f.getExportTask(ctx, appID, secret, tickit, sheetToken, usserAccessToken)
				if err != nil {
					f.logger.Error("get export task failed", log.Error(err))
					return nil
				}
				if res.FileName != "" && res.FileToken != "" {
					ticker.Stop()
					return res
				}
			case <-ctx.Done():
				ticker.Stop()
				return nil
			case <-time.After(3 * time.Second):
				ticker.Stop()
				f.logger.Error("get export task timeout", log.String("ticket", tickit))
				return nil
			}
		}
	})

	successResults := lo.Filter(getExportTaskResps, func(resp *getExportTaskResp, _ int) bool {
		return resp != nil
	})
	if len(successResults) < len(tickits) {
		return "", "", fmt.Errorf("some getExportTask failed (%d/%d succeeded)", len(successResults), len(tickits))
	}

	for i, resp := range successResults {
		f.logger.Debug("Download file info",
			log.Int("index", i),
			log.String("fileName", resp.FileName),
			log.String("fileToken", resp.FileToken))
	}

	res := parallel.Map(successResults, func(resp *getExportTaskResp, _ int) *domain.ScrapeResp {
		fileName, file, err := f.downloadExportTask(ctx, appID, secret, resp.FileToken, usserAccessToken)
		if err != nil {
			f.logger.Error("download export task failed", log.Error(err))
			return nil
		}

		imgName := fmt.Sprintf("%s/%s%s", kbID, uuid.New().String(), filepath.Ext(fileName))
		buf, _ := io.ReadAll(file)
		_, err = f.minioClient.PutObject(
			ctx,
			domain.Bucket,
			imgName,
			bytes.NewBuffer(buf),
			int64(len(buf)),
			minio.PutObjectOptions{
				ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				UserMetadata: map[string]string{
					"originalname": fileName,
				},
			},
		)
		if err != nil {
			f.logger.Error("upload export task failed", log.Error(err))
			return nil
		}
		ossPath := fmt.Sprintf("/%s/%s", domain.Bucket, imgName)
		res, err := f.crawler.ScrapeURL(ctx, ossPath, kbID)
		if err != nil {
			f.logger.Error("scrape export task failed", log.Error(err))
			return nil
		}
		return res
	})

	successResults2 := lo.Filter(res, func(resp *domain.ScrapeResp, _ int) bool {
		return resp != nil
	})
	for i, resp := range successResults2 {
		f.logger.Debug("Scraped content",
			log.Int("index", i),
			log.String("title", resp.Title),
			log.Int("content_length", len(resp.Content)))
	}
	if len(successResults2) < len(res) {
		return "", "", fmt.Errorf("some downloadExportTask failed (%d/%d succeeded)", len(successResults2), len(res))
	}
	info, err := getSheetInfo(ctx, appID, secret, sheetToken, usserAccessToken)
	if err != nil {
		f.logger.Error("get sheet info failed", log.Error(err))
		return "", "", err
	}
	// 在最后部分
	var (
		title   string
		content strings.Builder
	)

	// 添加调试信息
	f.logger.Debug("Processing results",
		log.Int("sheets_count", len(sheets)),
		log.Int("success_results_count", len(successResults2)))

	title = *info.Title
	if title == "" {
		title = "未命名表格"
	}
	content.WriteString(fmt.Sprintf("# %s\n\n", title))
	// content.WriteString(successResults2[0].Content)
	for _, successResult2 := range successResults2 {
		content.WriteString(successResult2.Content)
	}
	return title, content.String(), nil
}

type getExportTaskResp struct {
	FileName  string
	FileToken string
}

func (f *FeishuUseCase) getExportTask(ctx context.Context, appID, secret, ticket, sheetToken, userAccessToken string) (*getExportTaskResp, error) {
	client := lark.NewClient(appID, secret)
	req := larkdrive1.NewGetExportTaskReqBuilder().
		Ticket(ticket).
		Token(sheetToken).
		Build()

	// 发起请求
	resp, err := client.Drive.V1.ExportTask.Get(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
	if err != nil {
		return nil, fmt.Errorf("get export task failed: %v", err)
	}
	if !resp.Success() {
		return nil, fmt.Errorf("get export task failed: %s", resp.CodeError.Msg)
	}

	return &getExportTaskResp{
		FileName:  *resp.Data.Result.FileName,
		FileToken: *resp.Data.Result.FileToken,
	}, nil
}

func (f *FeishuUseCase) downloadExportTask(ctx context.Context, appID, secret, fileToken, userAccessToken string) (string, io.Reader, error) {
	client := lark.NewClient(appID, secret)
	// 创建请求对象
	req := larkdrive1.NewDownloadExportTaskReqBuilder().
		FileToken(fileToken).
		Build()

	// 发起请求
	resp, err := client.Drive.V1.ExportTask.Download(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
	// 处理错误
	if err != nil {
		return "", nil, fmt.Errorf("下载文件失败: %v", err)
	}
	return resp.FileName, resp.File, nil
}

func (f *FeishuUseCase) creatExportTask(ctx context.Context, appID, secret, token, userAccessToken string, subIDs []string) ([]string, error) {
	client := lark.NewClient(appID, secret)
	res := parallel.Map(subIDs, func(subID string, _ int) string {
		req := larkdrive1.NewCreateExportTaskReqBuilder().
			ExportTask(larkdrive1.NewExportTaskBuilder().
				FileExtension(`xlsx`).
				Token(token).
				Type(`sheet`).
				SubId(subID).
				Build()).
			Build()

		resp, err := client.Drive.V1.ExportTask.Create(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
		if err != nil {
			return ""
		}
		f.logger.Debug("Export task created", log.String("subId", subID), log.String("ticket", *resp.Data.Ticket))
		return *resp.Data.Ticket
	})
	// 处理错误
	successResults := lo.Filter(res, func(resp string, _ int) bool {
		return resp != ""
	})
	if len(successResults) < len(subIDs) {
		return nil, fmt.Errorf("some exports failed (%d/%d succeeded)", len(successResults), len(subIDs))
	}
	return res, nil
}

func (f *FeishuUseCase) getSheets(ctx context.Context, appID, secret, sheetToken, userAccessToken string) ([]*larksheets.Sheet, error) {
	// 创建 Client
	client := lark.NewClient(appID, secret)
	// 创建请求对象
	req := larksheets.NewQuerySpreadsheetSheetReqBuilder().
		SpreadsheetToken(sheetToken).
		Build()

	// 发起请求
	resp, err := client.Sheets.V3.SpreadsheetSheet.Query(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
	// 处理错误
	if err != nil {
		f.logger.Error("get sheets failed", log.Error(err))
		return nil, fmt.Errorf("get sheets failed: %v", err)
	}

	// 服务端错误处理
	if !resp.Success() {
		f.logger.Error("get sheets failed", log.String("logId", resp.RequestId()), log.String("error", larkcore.Prettify(resp.CodeError)))
		return nil, fmt.Errorf("get sheets failed: %s", resp.CodeError.Msg)
	}
	return resp.Data.Sheets, nil
}
