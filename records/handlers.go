package records

import (
	"amarkezic.github.com/finance-app/core"
)

var GetRecords = core.List[core.Record]()

var GetRecord = core.Single[core.Record]()

var CreateRecord = core.Create[core.Record]()

var UpdateRecord = core.Update[core.Record]()

var DeleteRecord = core.Delete[core.Record]()
