package http

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/itering/subscan/plugins"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/ss58"
)

func metadata(c *gin.Context) {
	metadata, err := svc.Metadata()
	util.ToJson(c, metadata, err)
}

func blocks(c *gin.Context) {
	p := new(struct {
		Row  int `json:"row" binding:"min=1,max=100"`
		Page int `json:"page" binding:"min=0"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}
	blockNum, err := svc.GetCurrentBlockNum(context.TODO())
	blocks := svc.GetBlocksSampleByNums(p.Page, p.Row)
	util.ToJson(c, map[string]interface{}{
		"blocks": blocks, "count": blockNum,
	}, err)
}

func block(c *gin.Context) {
	p := new(struct {
		BlockNum  int    `json:"block_num" binding:"omitempty,min=0"`
		BlockHash string `json:"block_hash" binding:"omitempty,len=66"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}
	if p.BlockHash == "" {
		util.ToJson(c, svc.GetBlockByNum(p.BlockNum), nil)
	} else {
		util.ToJson(c, svc.GetBlockByHashJson(p.BlockHash), nil)
	}
}

func extrinsics(c *gin.Context) {
	p := new(struct {
		Row     int    `json:"row" binding:"min=1,max=100"`
		Page    int    `json:"page" binding:"min=0"`
		Signed  string `json:"signed" binding:"omitempty"`
		Address string `json:"address" binding:"omitempty"`
		Module  string `json:"module" binding:"omitempty"`
		Call    string `json:"call" binding:"omitempty"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}
	var query []string
	if p.Module != "" {
		query = append(query, fmt.Sprintf("call_module = '%s'", p.Module))
	}
	if p.Call != "" {
		query = append(query, fmt.Sprintf("call_module_function = '%s'", p.Call))
	}

	if p.Signed == "signed" {
		query = append(query, "is_signed = 1")
	}
	if p.Address != "" {
		account := ss58.Decode(p.Address, util.StringToInt(util.AddressType))
		if account == "" {
			util.ToJson(c, nil, util.InvalidAccountAddress)
			return
		}
		query = append(query, fmt.Sprintf("is_signed = 1 and account_id = '%s'", account))
	}
	list, count := svc.GetExtrinsicList(p.Page, p.Row, "desc", query...)
	util.ToJson(c, map[string]interface{}{
		"extrinsics": list, "count": count,
	}, nil)
}

func extrinsic(c *gin.Context) {
	p := new(struct {
		ExtrinsicIndex string `json:"extrinsic_index" binding:"omitempty"`
		Hash           string `json:"hash" binding:"omitempty,len=66"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}
	if p.ExtrinsicIndex == "" && p.Hash == "" {
		util.ToJson(c, nil, util.ParamsError)
		return
	}
	if p.ExtrinsicIndex != "" {
		util.ToJson(c, svc.GetExtrinsicByIndex(p.ExtrinsicIndex), nil)
	} else {
		util.ToJson(c, svc.GetExtrinsicDetailByHash(p.Hash), nil)
	}
}

func events(c *gin.Context) {
	p := new(struct {
		Row    int    `json:"row" binding:"min=1,max=100"`
		Page   int    `json:"page" binding:"min=0"`
		Module string `json:"module" binding:"omitempty"`
		Call   string `json:"call" binding:"omitempty"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}
	var query []string
	if p.Module != "" {
		query = append(query, fmt.Sprintf("module_id = '%s'", p.Module))
	}
	if p.Call != "" {
		query = append(query, fmt.Sprintf("event_id = '%s'", p.Call))
	}
	events, count := svc.RenderEvents(p.Page, p.Row, "desc", query...)
	util.ToJson(c, map[string]interface{}{
		"events": events, "count": count,
	}, nil)
}

func checkSearchHash(c *gin.Context) {
	p := new(struct {
		Hash string `json:"hash" binding:"len=66"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}
	if block := svc.GetBlockByHash(p.Hash); block != nil {
		util.ToJson(c, map[string]string{"hash_type": "block"}, nil)
		return
	}
	if extrinsic := svc.GetExtrinsicByHash(p.Hash); extrinsic != nil {
		util.ToJson(c, map[string]string{"hash_type": "extrinsic"}, nil)
		return
	}
	util.ToJson(c, nil, util.RecordNotFound)
}

func runtimeList(c *gin.Context) {
	util.ToJson(c, map[string]interface{}{
		"list": svc.SubstrateRuntimeList(),
	}, nil)
}

func runtimeMetadata(c *gin.Context) {
	p := new(struct {
		Spec int `json:"spec"`
	})
	if err := c.MustBindWith(p, binding.JSON); err != nil {
		return
	}
	if info := svc.SubstrateRuntimeInfo(p.Spec); info == nil {
		util.ToJson(c, map[string]interface{}{"info": nil}, nil)
	} else {
		util.ToJson(c, map[string]interface{}{
			"info": info.Metadata.Modules,
		}, nil)
	}
}

func pluginList(c *gin.Context) {
	util.ToJson(c, plugins.List(), nil)
}
