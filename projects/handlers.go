package projects

import (
	"amarkezic.github.com/finance-app/core"
)

var GetProjects = core.List[core.Project]()

var GetProject = core.Single[core.Project]()

var CreateProject = core.Create[core.Project]()

var UpdateProject = core.Update[core.Project]()

var DeleteProject = core.Delete[core.Project]()
