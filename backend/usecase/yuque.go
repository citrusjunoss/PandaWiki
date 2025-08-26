package usecase

import (
	"archive/tar"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/utils"
	"gopkg.in/yaml.v2"
)

type YuqueUsecase struct {
	logger      *log.Logger
	fileusecase *FileUsecase
}

func NewYuqueUsecase(logger *log.Logger, fileusecase *FileUsecase) *YuqueUsecase {
	return &YuqueUsecase{
		logger:      logger.WithModule("usecase.yuqueusecase"),
		fileusecase: fileusecase,
	}
}

func (y *YuqueUsecase) AnalysisExportFile(ctx context.Context, fileHeader *multipart.FileHeader, kbID string) ([]domain.YuqueResp, error) {
	// 调用转换函数
	results, err := y.ConvertYuqueToMarkdown(fileHeader)
	if err != nil {
		y.logger.Error("Failed to convert Yuque to Markdown", "error", err)
		return nil, err
	}
	for i := range results {
		results[i].Content, err = y.exchangeUrl(ctx, results[i].Content, kbID)
		if err != nil {
			y.logger.Error("Failed to exchange Yuque URL", "error", err)
			return nil, err
		}
	}
	return results, nil
}

func (y *YuqueUsecase) exchangeUrl(ctx context.Context, md, kbid string) (string, error) {
	res, err := utils.ExchangeMarkDownImageUrl(ctx, []byte(md), func(ctx context.Context, originUrl *string) (string, error) {
		if originUrl == nil {
			y.logger.Error("Origin URL is nil")
			return "", nil
		}
		resp, err := http.Get(*originUrl)
		if err != nil {
			y.logger.Error("Failed to get image from URL", log.String("origin_url", *originUrl), log.String("error", err.Error()))
			return *originUrl, nil
		}
		defer resp.Body.Close()
		key, err := y.fileusecase.UploadFileFromReader(ctx, kbid, *originUrl, resp.Body, resp.ContentLength)
		return fmt.Sprintf("/%s/%s", domain.Bucket, key), err
	})
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// yuqueTocItem 表示目录项
type yuqueTocItem struct {
	Type  string `yaml:"type"`
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
	Level int    `yaml:"level"`
}

// yuqueDoc 结构体表示单个文档
type yuqueDoc struct {
	Body  string `json:"body"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

// yuqueDocument 结构体表示文档文件的结构
type yuqueDocument struct {
	Doc yuqueDoc `json:"doc"`
}

// yuqueMetaData 结构体表示元数据文件的结构
type yuqueMetaData struct {
	Book     map[string]interface{}   `json:"book"`
	Menus    []map[string]interface{} `json:"menus"`
	TocItems []yuqueTocItem           `json:"-"` // 不通过JSON解析，而是通过tocYml解析
}

// processDocument 处理单个文档
func processDocument(reader io.Reader, docLevelMap map[string]int) (*domain.YuqueResp, error) {
	// 读取文档文件内容
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// 解析JSON
	var document yuqueDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, err
	}

	mdContent, err := htmltomarkdown.ConvertString(document.Doc.Body)
	if err != nil {
		return nil, err
	}

	// 移除 <!--THE END--> 标记
	mdContent = strings.ReplaceAll(mdContent, "<!--THE END-->", "")

	// 创建返回结果
	resp := &domain.YuqueResp{
		Content: mdContent,
		Title:   document.Doc.Title,
	}

	return resp, nil
}

// processMeta 处理元数据文件
func processMeta(reader io.Reader) (*yuqueMetaData, error) {
	// 读取元数据文件内容
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// 解析JSON
	var metaData yuqueMetaData
	if err := json.Unmarshal(data, &metaData); err != nil {
		return nil, err
	}

	// 解析tocYml
	if book, ok := metaData.Book["tocYml"].(string); ok {
		var tocItems []yuqueTocItem
		if err := yaml.Unmarshal([]byte(book), &tocItems); err != nil {
			return nil, fmt.Errorf("failed to parse tocYml: %w", err)
		}
		metaData.TocItems = tocItems
	}

	return &metaData, nil
}

// ConvertYuqueToMarkdown 将语雀文件转换为Markdown
func (y *YuqueUsecase) ConvertYuqueToMarkdown(fileHeader *multipart.FileHeader) ([]domain.YuqueResp, error) {
	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		y.logger.Error("Failed to open file header", "error", err)
		return nil, err
	}
	defer file.Close()

	// 创建tar reader
	tarReader := tar.NewReader(file)

	// 用于存储元数据
	var metaData *yuqueMetaData
	// 用于存储文档URL到层级的映射
	docLevelMap := make(map[string]int)
	// 用于存储结果
	var results []domain.YuqueResp

	// 遍历tar中的每个文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			y.logger.Error("Failed to read tar header", "error", err)
			return nil, err
		}
		fileName := utils.RemoveFirstDir(header.Name)
		// 处理元数据文件
		if fileName == "$meta.json" {
			metaData, err = processMeta(tarReader)
			if err != nil {
				return nil, err
			}

			// 创建一个映射来存储文档URL到层级的映射
			for _, item := range metaData.TocItems {
				if item.Type == "DOC" {
					docLevelMap[item.URL] = item.Level
				}
			}
			continue
		}

		// 处理文档文件
		if strings.HasSuffix(header.Name, ".json") {
			resp, err := processDocument(tarReader, docLevelMap)
			if err != nil {
				return nil, err
			}
			results = append(results, *resp)
		}
	}

	// 检查是否找到了元数据文件
	if metaData == nil {
		return nil, fmt.Errorf("meta data not found in tar file")
	}

	return results, nil
}
