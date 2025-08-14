package dangbei

import (
	"context"
	"time"

	"github.com/WenChunTech/OpenAICompatible/src/model/openai"
)

//	{
//	    "success": true,
//	    "errCode": null,
//	    "errMessage": null,
//	    "requestId": "8aea1d64-d215-427e-8ee6-9d32bb6aa01d",
//	    "data": {
//	        "model": "deepseek",
//	        "modelList": [
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/db-ai-search/2025/07/01/newbanner%402x.png",
//	                "bannerBadge": "https://ai-search-static.dangbei.net/2025/06/18/1935219346766958592.png",
//	                "icon": "https://ai-search-static.dangbei.net/2025/03/21/1902984234851766272.png",
//	                "title": "DeepSeek-R1最新版",
//	                "value": "deepseek",
//	                "badge": "https://ai-search-static.dangbei.net/2025/04/03/1907633212834844672.png",
//	                "innerBadgeText": "HOT",
//	                "hoverTitle": "DeepSeek-R1最新版",
//	                "hoverText": "专注逻辑推理与深度分析，擅长解决复杂问题，提供精准决策支持",
//	                "option": [
//	                    {
//	                        "disable": false,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": true
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": true,
//	                "recently": false,
//	                "pinned": true
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/2025/06/18/1935214661066690560.png",
//	                "bannerBadge": null,
//	                "icon": "https://ai-search-static.dangbei.net/2025/03/21/1902984234419752960.png",
//	                "title": "豆包-1.6",
//	                "value": "doubao-1_6-thinking",
//	                "badge": null,
//	                "innerBadgeText": null,
//	                "hoverTitle": "豆包（doubao-1.6-thinking）",
//	                "hoverText": "豆包最新推理模型，创作、推理、数学大幅增强",
//	                "option": [
//	                    {
//	                        "disable": false,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": true
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": true,
//	                "recently": false,
//	                "pinned": true
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/2025/06/18/1935214664736706560.png",
//	                "bannerBadge": null,
//	                "icon": "https://ai-search-static.dangbei.net/2025/03/21/1902984234851766272.png",
//	                "title": "DeepSeek-V3",
//	                "value": "deepseek-v3",
//	                "badge": null,
//	                "innerBadgeText": null,
//	                "hoverTitle": "DeepSeek-V3",
//	                "hoverText": "轻量高效，响应极快。擅长代码，可高效解析代码与图表",
//	                "option": [
//	                    {
//	                        "disable": true,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": false
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": false,
//	                "recently": false,
//	                "pinned": false
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/db-ai-search/2025/07/31/zhipu_glm_4_5_banner.png",
//	                "bannerBadge": null,
//	                "icon": "https://ai-search-static.dangbei.net/db-ai-search/2025/07/21/zhipu_logo.png",
//	                "title": "GLM-4.5",
//	                "value": "glm-4-5",
//	                "badge": "https://ai-search-static.dangbei.net/2025/04/03/2a55c10c-00f5-42e9-867c-ca5e43251e78.png",
//	                "innerBadgeText": "NEW",
//	                "hoverTitle": "GLM-4.5",
//	                "hoverText": "智谱最新旗舰模型，支持思考模式切换，综合能力达到开源模型的SOTA水平。",
//	                "option": [
//	                    {
//	                        "disable": false,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": false
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": false,
//	                "recently": false,
//	                "pinned": false
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/db-ai-search/2025/07/21/qwen3-235b-banner.png",
//	                "bannerBadge": "https://ai-search-static.dangbei.net/2025/06/18/1935219346766958592.png",
//	                "icon": "https://ai-search-static.dangbei.net/2025/03/21/1902984233610252288.png",
//	                "title": "通义3-235B",
//	                "value": "qwen3-235b-a22b",
//	                "badge": "https://ai-search-static.dangbei.net/2025/04/03/2a55c10c-00f5-42e9-867c-ca5e43251e78.png",
//	                "innerBadgeText": "0722",
//	                "hoverTitle": "通义千问（qwen3-235b）",
//	                "hoverText": "国内首个混合推理模型，达到同规模业界SOTA水平。",
//	                "option": [
//	                    {
//	                        "disable": false,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": false
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": false,
//	                "recently": false,
//	                "pinned": false
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/db-ai-search/2025/07/15/1944986666569699328/kimi-k2.png",
//	                "bannerBadge": "https://ai-search-static.dangbei.net/2025/06/18/1935219346766958592.png",
//	                "icon": "https://ai-search-static.dangbei.net/2025/03/30/1902984234004516864.png",
//	                "title": "Kimi K2",
//	                "value": "kimi-k2-0711-preview",
//	                "badge": "https://ai-search-static.dangbei.net/2025/04/03/2a55c10c-00f5-42e9-867c-ca5e43251e78.png",
//	                "innerBadgeText": "NEW",
//	                "hoverTitle": "Kimi K2（kimi-k2-0711）",
//	                "hoverText": "具备更强代码能力、更擅长通用Agent任务",
//	                "option": [
//	                    {
//	                        "disable": true,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": false
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": false,
//	                "recently": false,
//	                "pinned": false
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/db-ai-search/2025/07/21/minimax_banner_0721.png",
//	                "bannerBadge": null,
//	                "icon": "https://ai-search-static.dangbei.net/db-ai-search/2025/07/21/minimax_logo.png",
//	                "title": "MiniMax-M1",
//	                "value": "MiniMax-M1",
//	                "badge": null,
//	                "innerBadgeText": null,
//	                "hoverTitle": "MiniMax-M1",
//	                "hoverText": "全球领先，80K思维链 x 1M输入",
//	                "option": [
//	                    {
//	                        "disable": false,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": false
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": false,
//	                "recently": false,
//	                "pinned": false
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/db-ai-search/2025/07/21/zhipu_banner_0721.png",
//	                "bannerBadge": null,
//	                "icon": "https://ai-search-static.dangbei.net/db-ai-search/2025/07/21/zhipu_logo.png",
//	                "title": "GLM-4-Plus",
//	                "value": "glm-4-plus",
//	                "badge": null,
//	                "innerBadgeText": null,
//	                "hoverTitle": "智谱清言（GLM-4-Plus）",
//	                "hoverText": "智谱最强高智能旗舰模型",
//	                "option": [
//	                    {
//	                        "disable": true,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": false
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": false,
//	                "recently": false,
//	                "pinned": false
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/2025/06/18/1935214661708419072.png",
//	                "bannerBadge": "https://ai-search-static.dangbei.net/2025/06/18/1935219346766958592.png",
//	                "icon": "https://ai-search-static.dangbei.net/2025/03/21/1902984234419752960.png",
//	                "title": "豆包",
//	                "value": "doubao",
//	                "badge": null,
//	                "innerBadgeText": null,
//	                "hoverTitle": "豆包（doubao-1.5-pro-32k）",
//	                "hoverText": "字节全能AI，创意写作、百科解答、难题破解，随需响应",
//	                "option": [
//	                    {
//	                        "disable": true,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": false
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": false,
//	                "recently": false,
//	                "pinned": false
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/2025/06/18/1935214660353658880.png",
//	                "bannerBadge": "https://ai-search-static.dangbei.net/2025/06/18/1935219346766958592.png",
//	                "icon": "https://ai-search-static.dangbei.net/2025/03/21/1902984233610252288.png",
//	                "title": "通义Plus",
//	                "value": "qwen-plus",
//	                "badge": null,
//	                "innerBadgeText": null,
//	                "hoverTitle": "通义千问（qwen-plus-32k）",
//	                "hoverText": "复杂问题速解专家，知识广博，表达清晰精准",
//	                "option": [
//	                    {
//	                        "disable": true,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": false
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": false,
//	                "recently": false,
//	                "pinned": false
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/2025/06/18/1935214664422133760.png",
//	                "bannerBadge": "https://ai-search-static.dangbei.net/2025/06/18/1935219346766958592.png",
//	                "icon": "https://ai-search-static.dangbei.net/2025/03/30/1902984234004516864.png",
//	                "title": "Kimi",
//	                "value": "moonshot-v1-32k",
//	                "badge": null,
//	                "innerBadgeText": null,
//	                "hoverTitle": "Kimi（moonshot-v1-32k）",
//	                "hoverText": "高效问题解析者，多领域知识库，语言简练有力",
//	                "option": [
//	                    {
//	                        "disable": true,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": false
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": false,
//	                "recently": false,
//	                "pinned": false
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/2025/06/18/1935214659141505024.png",
//	                "bannerBadge": "https://ai-search-static.dangbei.net/2025/06/18/1935219346766958592.png",
//	                "icon": "https://ai-search-static.dangbei.net/2025/03/21/1902984233610252288.png",
//	                "title": "通义QwQ",
//	                "value": "qwq-plus",
//	                "badge": null,
//	                "innerBadgeText": null,
//	                "hoverTitle": "通义千问（qwq-plus）",
//	                "hoverText": "善解难题，精准表达，知识全面",
//	                "option": [
//	                    {
//	                        "disable": false,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": true
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": false,
//	                "recently": false,
//	                "pinned": false
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/2025/06/18/1935214660760506368.png",
//	                "bannerBadge": "https://ai-search-static.dangbei.net/2025/06/18/1935219346766958592.png",
//	                "icon": "https://ai-search-static.dangbei.net/2025/03/21/1902984233610252288.png",
//	                "title": "通义Long",
//	                "value": "qwen-long",
//	                "badge": null,
//	                "innerBadgeText": null,
//	                "hoverTitle": "通义千问（qwen-long）",
//	                "hoverText": "通义千问针对超长上下文处理场景的大语言模型",
//	                "option": [
//	                    {
//	                        "disable": true,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": false
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": false,
//	                "recently": false,
//	                "pinned": false
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/2025/06/18/1935214661364486144.png",
//	                "bannerBadge": "https://ai-search-static.dangbei.net/2025/06/18/1935219346766958592.png",
//	                "icon": "https://ai-search-static.dangbei.net/2025/03/21/1902984234419752960.png",
//	                "title": "豆包-1.5",
//	                "value": "doubao-thinking",
//	                "badge": null,
//	                "innerBadgeText": null,
//	                "hoverTitle": "豆包（doubao-1.5-thinking-pro）",
//	                "hoverText": "推理模型，专精数理编程，擅长创意写作",
//	                "option": [
//	                    {
//	                        "disable": false,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": true
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": false,
//	                "recently": false,
//	                "pinned": false
//	            },
//	            {
//	                "banner": "https://ai-search-static.dangbei.net/2025/06/18/1935214664132726784.png",
//	                "bannerBadge": "https://ai-search-static.dangbei.net/2025/06/18/1935219346766958592.png",
//	                "icon": "https://ai-search-static.dangbei.net/2025/04/15/wenxin_icon.png",
//	                "title": "文心4.5",
//	                "value": "ernie-4.5-turbo-32k",
//	                "badge": null,
//	                "innerBadgeText": null,
//	                "hoverTitle": "文心一言（ernie-4.5-turbo-32k）",
//	                "hoverText": "广泛适用于各领域复杂任务场景",
//	                "option": [
//	                    {
//	                        "disable": true,
//	                        "title": "深度思考",
//	                        "value": "deep",
//	                        "selected": false
//	                    },
//	                    {
//	                        "disable": false,
//	                        "title": "联网搜索",
//	                        "value": "online",
//	                        "selected": false
//	                    }
//	                ],
//	                "recommend": false,
//	                "recently": false,
//	                "pinned": false
//	            }
//	        ]
//	    }
//	}
type DangBeiModelListResponse struct {
	Success    bool   `json:"success"`
	ErrCode    any    `json:"errCode"`
	ErrMessage any    `json:"errMessage"`
	RequestID  string `json:"requestId"`
	Data       Data   `json:"data"`
}
type Option struct {
	Disable  bool   `json:"disable"`
	Title    string `json:"title"`
	Value    string `json:"value"`
	Selected bool   `json:"selected"`
}
type ModelList struct {
	Banner         string   `json:"banner"`
	BannerBadge    string   `json:"bannerBadge"`
	Icon           string   `json:"icon"`
	Title          string   `json:"title"`
	Value          string   `json:"value"`
	Badge          string   `json:"badge"`
	InnerBadgeText string   `json:"innerBadgeText"`
	HoverTitle     string   `json:"hoverTitle"`
	HoverText      string   `json:"hoverText"`
	Option         []Option `json:"option"`
	Recommend      bool     `json:"recommend"`
	Recently       bool     `json:"recently"`
	Pinned         bool     `json:"pinned"`
}
type Data struct {
	Model     string      `json:"model"`
	ModelList []ModelList `json:"modelList"`
}

func (c *DangBeiModelListResponse) Convert(ctx context.Context) (*openai.OpenAIModelListResponse, error) {
	models := make([]*openai.Model, len(c.Data.ModelList))
	for i, model := range c.Data.ModelList {
		models[i] = &openai.Model{
			ID:      model.Value,
			Object:  "model",
			Created: time.Now().Unix(),
			OwnedBy: "dangbei",
		}
	}

	return &openai.OpenAIModelListResponse{
		Object: "list",
		Data:   models,
	}, nil
}
