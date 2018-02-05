# go-harvest #

go-harvest is a Go client library for accessing [Harvest API v2](https://help.getharvest.com/api-v2/)

## Usage ##

```go
import "github.com/becoded/go-harvest/harvest"
```

Construct a new Harvest client, then use the various services on the client to
access different parts of the Harvest API. For example:
```
ctx := context.Background()
ts := oauth2.StaticTokenSource(
	&oauth2.Token{
        AccessToken: os.Getenv("HARVEST_ACCESS_TOKEN"),
    },
)
tc := oauth2.NewClient(ctx, ts)
service := harvest.NewHarvestClient(tc)
service.AccountId = os.Getenv("HARVEST_ACCOUNT_ID")

c, _, err := service.Company.Get(ctx)
if err != nil {
    fmt.Print(err)
    panic("Exit company")
}
```


##[API Introduction](https://help.getharvest.com/api-v2/introduction)
* [Overview](https://help.getharvest.com/api-v2/introduction/overview/general/)
* [Code Samples](https://help.getharvest.com/api-v2/introduction/overview/code-samples/)
* [Supported Time Zones](https://help.getharvest.com/api-v2/introduction/overview/supported-timezones/)
* [Supported Currencies](https://help.getharvest.com/api-v2/introduction/overview/supported-currencies/)
* [Pagination](https://help.getharvest.com/api-v2/introduction/overview/pagination/)

###Legend
* [ - ] Partially done
* [ x ] Complete


##[Authentication](https://help.getharvest.com/api-v2/authentication-api)
* [ - ] [Authentication](https://help.getharvest.com/api-v2/authentication-api/authentication/authentication/)

##[Clients API](https://help.getharvest.com/api-v2/clients-api)
* [ x ] [Client Contacts](https://help.getharvest.com/api-v2/clients-api/clients/contacts/)
* [ x ] [Clients](https://help.getharvest.com/api-v2/clients-api/clients/clients/)

##[Company Settings](https://help.getharvest.com/api-v2/company-api)
* [ x ] [Company](https://help.getharvest.com/api-v2/company-api/company/company/)

##[Invoices API](https://help.getharvest.com/api-v2/invoices-api)
* [ x ] [Invoice Messages](https://help.getharvest.com/api-v2/invoices-api/invoices/invoice-messages/)
* [ x ] [Invoice Payments](https://help.getharvest.com/api-v2/invoices-api/invoices/invoice-payments/)
* [ - ] [Invoices](https://help.getharvest.com/api-v2/invoices-api/invoices/invoices/)
* [ - ] [Invoice Item Categories](https://help.getharvest.com/api-v2/invoices-api/invoices/invoice-item-categories/)

##[Estimates API](https://help.getharvest.com/api-v2/estimates-api)
* [ - ] [Estimate Messages](https://help.getharvest.com/api-v2/estimates-api/estimates/estimate-messages/)
* [ - ] [Estimates](https://help.getharvest.com/api-v2/estimates-api/estimates/estimates/)
* [ - ] [Estimate Item Categories](https://help.getharvest.com/api-v2/estimates-api/estimates/estimate-item-categories/)

##[Expenses API](https://help.getharvest.com/api-v2/expenses-api)
* [ - ] [Expenses](https://help.getharvest.com/api-v2/expenses-api/expenses/expenses/)
* [ - ] [Expense Categories](https://help.getharvest.com/api-v2/expenses-api/expenses/expense-categories/)

##[Tasks API](https://help.getharvest.com/api-v2/tasks-api)
* [ - ] [Tasks](https://help.getharvest.com/api-v2/tasks-api/tasks/tasks/)

##[Timesheets API](https://help.getharvest.com/api-v2/timesheets-api)
* [ - ] [Time Entries](https://help.getharvest.com/api-v2/timesheets-api/timesheets/time-entries/)

##[Projects API](https://help.getharvest.com/api-v2/projects-api)
* [ - ] [Project User Assignments](https://help.getharvest.com/api-v2/projects-api/projects/user-assignments/)
* [ - ] [Project Task Assignments](https://help.getharvest.com/api-v2/projects-api/projects/task-assignments/)
* [ - ] [Projects](https://help.getharvest.com/api-v2/projects-api/projects/projects/)

##[Roles API](https://help.getharvest.com/api-v2/roles-api)
* [ X ] [Roles](https://help.getharvest.com/api-v2/roles-api/roles/roles/)

##[Users API](https://help.getharvest.com/api-v2/users-api)
* [ - ] [User Project Assignments](https://help.getharvest.com/api-v2/users-api/users/project-assignments/)
* [ - ] [Users](https://help.getharvest.com/api-v2/users-api/users/users/)

## Todo
* [ ] Unit tests
* [ ] Rate limits
* [ ] Documentation
* [ ] Examples
* [ ] License