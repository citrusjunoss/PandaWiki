package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/repo/cache"
	"github.com/chaitin/panda-wiki/repo/ipdb"
	"github.com/chaitin/panda-wiki/repo/pg"
)

type StatUseCase struct {
	repo             *pg.StatRepository
	nodeRepo         *pg.NodeRepository
	conversationRepo *pg.ConversationRepository
	appRepo          *pg.AppRepository
	ipRepo           *ipdb.IPAddressRepo
	logger           *log.Logger
	geoCacheRepo     *cache.GeoRepo
	authRepo         *pg.AuthRepo
}

func NewStatUseCase(repo *pg.StatRepository, nodeRepo *pg.NodeRepository, conversationRepo *pg.ConversationRepository, appRepo *pg.AppRepository, ipRepo *ipdb.IPAddressRepo, geoCacheRepo *cache.GeoRepo, authRepo *pg.AuthRepo, logger *log.Logger) *StatUseCase {
	return &StatUseCase{
		repo:             repo,
		nodeRepo:         nodeRepo,
		conversationRepo: conversationRepo,
		appRepo:          appRepo,
		ipRepo:           ipRepo,
		geoCacheRepo:     geoCacheRepo,
		authRepo:         authRepo,
		logger:           logger.WithModule("usecase.stats"),
	}
}

func (u *StatUseCase) RecordPage(ctx context.Context, stat *domain.StatPage) error {
	if err := u.repo.CreateStatPage(ctx, stat); err != nil {
		return err
	}
	remoteIP := stat.IP
	ipAddress, err := u.ipRepo.GetIPAddress(ctx, remoteIP)
	if err != nil {
		u.logger.Warn("get ip address failed", log.Error(err), log.String("ip", remoteIP), log.Int64("stat_id", stat.ID))
	} else {
		location := fmt.Sprintf("%s|%s|%s", ipAddress.Country, ipAddress.Province, ipAddress.City)
		if err := u.geoCacheRepo.SetGeo(ctx, stat.KBID, location); err != nil {
			u.logger.Warn("set geo cache failed", log.Error(err), log.Int64("stat_id", stat.ID), log.String("ip", remoteIP))
		}
	}
	return nil
}

func (u *StatUseCase) GetHotPages(ctx context.Context, kbID string) ([]*domain.HotPageResp, error) {
	hotPages, err := u.repo.GetHotPages(ctx, kbID)
	if err != nil {
		return nil, err
	}
	nodeIDs := lo.Uniq(lo.Map(hotPages, func(page *domain.HotPageResp, _ int) string {
		return page.NodeID
	}))
	docNames, err := u.nodeRepo.GetNodeNameByNodeIDs(ctx, nodeIDs)
	if err != nil {
		return nil, err
	}
	for _, page := range hotPages {
		switch page.Scene {
		case domain.StatPageSceneNodeDetail:
			page.NodeName = docNames[page.NodeID]
		case domain.StatPageSceneWelcome:
			page.NodeName = "欢迎页"
		case domain.StatPageSceneChat:
			page.NodeName = "问答页"
		case domain.StatPageSceneLogin:
			page.NodeName = "登录页"
		}
	}
	return hotPages, nil
}

func (u *StatUseCase) GetHotRefererHosts(ctx context.Context, kbID string) ([]*domain.HotRefererHostResp, error) {
	hotRefererHosts, err := u.repo.GetHotRefererHosts(ctx, kbID)
	if err != nil {
		return nil, err
	}
	return hotRefererHosts, nil
}

func (u *StatUseCase) GetHotBrowsers(ctx context.Context, kbID string) (*domain.HotBrowserResp, error) {
	hotBrowsers, err := u.repo.GetHotBrowsers(ctx, kbID)
	if err != nil {
		return nil, err
	}
	return hotBrowsers, nil
}

func (u *StatUseCase) GetCount(ctx context.Context, kbID string) (*domain.StatPageCountResp, error) {
	count, err := u.repo.GetCount(ctx, kbID)
	if err != nil {
		return nil, err
	}
	conversationCount, err := u.conversationRepo.GetConversationCount(ctx, kbID)
	if err != nil {
		return nil, err
	}
	count.ConversationCount = conversationCount
	return count, nil
}

func (u *StatUseCase) GetInstantCount(ctx context.Context, kbID string) ([]*domain.InstantCountResp, error) {
	instantCount, err := u.repo.GetInstantCount(ctx, kbID)
	if err != nil {
		return nil, err
	}
	return instantCount, nil
}

func (u *StatUseCase) GetInstantPages(ctx context.Context, kbID string) ([]*domain.InstantPageResp, error) {
	pages, err := u.repo.GetInstantPages(ctx, kbID)
	if err != nil {
		return nil, err
	}
	ips := lo.Map(pages, func(page *domain.InstantPageResp, _ int) string {
		return page.IP
	})
	ipAddresses, err := u.ipRepo.GetIPAddresses(ctx, ips)
	if err != nil {
		return nil, err
	}
	authIDs := make([]uint, 0, 10)
	for _, page := range pages {
		ipAddress, ok := ipAddresses[page.IP]
		if !ok {
			ipAddress = &domain.IPAddress{
				IP:       page.IP,
				Country:  "未知",
				Province: "未知",
				City:     "未知",
			}
		}
		page.IPAddress = *ipAddress
		if page.UserID != 0 {
			authIDs = append(authIDs, page.UserID)
		}
	}
	authMap, err := u.authRepo.GetAuthUserinfoByIDs(ctx, authIDs)
	if err != nil {
		u.logger.Error("get user info failed", log.Error(err))
	}
	nodeIDs := lo.Uniq(lo.Map(pages, func(page *domain.InstantPageResp, _ int) string {
		return page.NodeID
	}))
	docNames, err := u.nodeRepo.GetNodeNameByNodeIDs(ctx, nodeIDs)
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		switch page.Scene {
		case domain.StatPageSceneNodeDetail:
			page.NodeName = docNames[page.NodeID]
		case domain.StatPageSceneWelcome:
			page.NodeName = "欢迎页"
		case domain.StatPageSceneChat:
			page.NodeName = "问答页"
		case domain.StatPageSceneLogin:
			page.NodeName = "登录页"
		default:
			page.NodeName = "未知"
		}
		if _, ok := authMap[page.UserID]; ok {
			page.Info = &domain.AuthUserInfo{
				Username:  authMap[page.UserID].AuthUserInfo.Username,
				Email:     authMap[page.UserID].AuthUserInfo.Email,
				AvatarUrl: authMap[page.UserID].AuthUserInfo.AvatarUrl,
			}
		}
	}
	return pages, nil
}

func (u *StatUseCase) GetGeoCount(ctx context.Context, kbID string) (map[string]int64, error) {
	geoCount, err := u.geoCacheRepo.GetLast24HourGeo(ctx, kbID)
	if err != nil {
		return nil, err
	}
	return geoCount, nil
}

func (u *StatUseCase) GetConversationDistribution(ctx context.Context, kbID string) ([]*domain.ConversationDistributionResp, error) {
	distribution, err := u.conversationRepo.GetConversationDistribution(ctx, kbID)
	if err != nil {
		return nil, err
	}
	// assign app_type
	appMap, err := u.appRepo.GetAppList(ctx, kbID)
	if err != nil {
		return nil, err
	}
	distribution = lo.Map(distribution, func(resp *domain.ConversationDistributionResp, _ int) *domain.ConversationDistributionResp {
		if dist, ok := appMap[resp.AppID]; ok {
			resp.AppType = dist.Type
		}
		return resp
	})
	return distribution, nil
}
