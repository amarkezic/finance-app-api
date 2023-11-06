package company

import (
	"amarkezic.github.com/finance-app/core"
)

var GetCompanies = core.List[core.Company]()

var GetCompany = core.Single[core.Company]()

var CreateCompany = core.Create[core.Company]()

var UpdateCompany = core.Update[core.Company]()

var DeleteCompany = core.Delete[core.Company]()
