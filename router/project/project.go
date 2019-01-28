// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package project

import (
    "github.com/gin-gonic/gin"
    "github.com/dreamans/syncd/render"
    "github.com/dreamans/syncd/module/project"
    "github.com/dreamans/syncd/util/gostring"
    "github.com/dreamans/syncd/util/goslice"
)

type ProjectFormBind struct {
    ID                  int     `form:"id"`
    SpaceId             int     `form:"space_id"`
    Name                string  `form:"name" binding:"required"`
    Description         string  `form:"description"`
    NeedAudit           int     `form:"need_audit"`
    RepoUrl             string  `form:"repo_url" binding:"required"`
    RepoBranch          string  `form:"repo_branch"`
    PreReleaseCluster   int     `form:"pre_release_cluster"`
    OnlineCluster       []int   `form:"online_cluster" binding:"required"`
    DeployUser          string  `form:"deploy_user" binding:"required"`
    DeployPath          string  `form:"deploy_path" binding:"required"`
    PreDeployCmd        string  `form:"pre_deploy_cmd"`
    AfterDeployCmd      string  `form:"after_deploy_cmd"`
    DeployTimeout       int     `form:"deploy_timeout" binding:"required"`
}

type QueryBind struct {
    SpaceId     int     `form:"space_id"`
    Keyword	    string  `form:"keyword"`
    Offset	    int     `form:"offset"`
    Limit	    int     `form:"limit" binding:"required,gte=1,lte=999"`
}

func ProjectDelete(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    proj := &project.Project{
        ID: id,
    }
    if err := proj.Delete(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.Success(c)
}

func ProjectDetail(c *gin.Context) {
    id := gostring.Str2Int(c.Query("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    proj := &project.Project{
        ID: id,
    }
    if err := proj.Detail(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.JSON(c, proj)
}

func ProjectSwitchStatus(c *gin.Context) {
    id, status := gostring.Str2Int(c.PostForm("id")), gostring.Str2Int(c.PostForm("status"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    if status !=0 {
        status = 1
    }
    proj := &project.Project{
        ID: id,
        Status: status,
    }
    if err := proj.UpdateStatus(); err != nil {
        render.AppError(c, err.Error())
        return
    }
    render.Success(c)
}

func ProjectList(c *gin.Context) {
    var query QueryBind
    if err := c.ShouldBind(&query); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    if query.SpaceId == 0 {
        render.ParamError(c, "space_id cannot be empty")
        return
    }
    proj := &project.Project{}
    list, err := proj.List(query.Keyword, query.SpaceId, query.Offset, query.Limit)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    total, err := proj.Total(query.Keyword, query.SpaceId)
    if err != nil {
        render.AppError(c, err.Error())
        return
    }

    projList := []map[string]interface{}{}
    for _, l := range list {
        projList = append(projList, map[string]interface{}{
            "id": l.ID,
            "name": l.Name,
            "need_audit": l.NeedAudit,
            "status": l.Status,
        })
    }

    render.JSON(c, gin.H{
        "list": projList,
        "total": total,
    })
}

func ProjectAdd(c *gin.Context) {
    projectCreateOrUpdate(c)
}

func ProjectUpdate(c *gin.Context) {
    id := gostring.Str2Int(c.PostForm("id"))
    if id == 0 {
        render.ParamError(c, "id cannot be empty")
        return
    }
    projectCreateOrUpdate(c)
}

func projectCreateOrUpdate(c *gin.Context) {
    var projectForm ProjectFormBind
    if err := c.ShouldBind(&projectForm); err != nil {
        render.ParamError(c, err.Error())
        return
    }
    onlineCluster := goslice.FilterSliceInt(projectForm.OnlineCluster)
    if len(onlineCluster) == 0 {
        render.ParamError(c, "online_cluster cannot be empty")
        return
    }
    proj := &project.Project{
        ID: projectForm.ID,
        SpaceId: projectForm.SpaceId,
        Name: projectForm.Name,
        Description: projectForm.Description,
        NeedAudit: projectForm.NeedAudit,
        RepoUrl: projectForm.RepoUrl,
        RepoBranch: projectForm.RepoBranch,
        PreReleaseCluster: projectForm.PreReleaseCluster,
        OnlineCluster: onlineCluster,
        DeployUser: projectForm.DeployUser,
        DeployPath: projectForm.DeployPath,
        PreDeployCmd: projectForm.PreDeployCmd,
        AfterDeployCmd: projectForm.AfterDeployCmd,
        DeployTimeout: projectForm.DeployTimeout,
    }
    if err := proj.CreateOrUpdate(); err != nil {
        render.AppError(c, err.Error())
        return 
    }
    render.Success(c)
}