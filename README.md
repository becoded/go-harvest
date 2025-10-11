# go-harvest #

[![GoDoc](https://godoc.org/github.com/becoded/go-harvest/harvest?status.svg)](https://godoc.org/github.com/becoded/go-harvest/harvest)
[![GitHub CI Test](https://github.com/becoded/go-harvest/actions/workflows/test.yml/badge.svg)](https://github.com/becoded/go-harvest/actions/workflows/test.yml)
[![GitHub CI Lint](https://github.com/becoded/go-harvest/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/becoded/go-harvest/actions/workflows/golangci-lint.yml)
[![Test Coverage](https://codecov.io/gh/becoded/go-harvest/branch/main/graph/badge.svg)](https://codecov.io/gh/becoded/go-harvest)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/3255/badge)](https://bestpractices.coreinfrastructure.org/projects/3255)
[![Go Report Card](https://goreportcard.com/badge/github.com/becoded/go-harvest)](https://goreportcard.com/report/github.com/becoded/go-harvest)

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
service := harvest.NewAPIClient(tc)
service.AccountID = os.Getenv("HARVEST_ACCOUNT_ID")

c, _, err := service.Company.Get(ctx)
if err != nil {
    fmt.Print(err)
    panic("Exit company")
}
```


## [API Introduction](https://help.getharvest.com/api-v2/introduction)
* [Overview](https://help.getharvest.com/api-v2/introduction/overview/general/)
* [Code Samples](https://help.getharvest.com/api-v2/introduction/overview/code-samples/)
* [Supported Time Zones](https://help.getharvest.com/api-v2/introduction/overview/supported-timezones/)
* [Supported Currencies](https://help.getharvest.com/api-v2/introduction/overview/supported-currencies/)
* [Pagination](https://help.getharvest.com/api-v2/introduction/overview/pagination/)

## [Authentication](https://help.getharvest.com/api-v2/authentication-api)
* [Authentication](https://help.getharvest.com/api-v2/authentication-api/authentication/authentication/)

## [Clients API](https://help.getharvest.com/api-v2/clients-api)
* [Client Contacts](https://help.getharvest.com/api-v2/clients-api/clients/contacts/)
* [Clients](https://help.getharvest.com/api-v2/clients-api/clients/clients/)

## [Company Settings](https://help.getharvest.com/api-v2/company-api)
* [Company](https://help.getharvest.com/api-v2/company-api/company/company/)

## [Invoices API](https://help.getharvest.com/api-v2/invoices-api)
* [Invoice Messages](https://help.getharvest.com/api-v2/invoices-api/invoices/invoice-messages/)
* [Invoice Payments](https://help.getharvest.com/api-v2/invoices-api/invoices/invoice-payments/)
* [Invoices](https://help.getharvest.com/api-v2/invoices-api/invoices/invoices/)
* [Invoice Item Categories](https://help.getharvest.com/api-v2/invoices-api/invoices/invoice-item-categories/)

## [Estimates API](https://help.getharvest.com/api-v2/estimates-api)
* [Estimate Messages](https://help.getharvest.com/api-v2/estimates-api/estimates/estimate-messages/)
* [Estimates](https://help.getharvest.com/api-v2/estimates-api/estimates/estimates/)
* [Estimate Item Categories](https://help.getharvest.com/api-v2/estimates-api/estimates/estimate-item-categories/)

## [Expenses API](https://help.getharvest.com/api-v2/expenses-api)
* [Expenses](https://help.getharvest.com/api-v2/expenses-api/expenses/expenses/)
* [Expense Categories](https://help.getharvest.com/api-v2/expenses-api/expenses/expense-categories/)

## [Tasks API](https://help.getharvest.com/api-v2/tasks-api)
* [Tasks](https://help.getharvest.com/api-v2/tasks-api/tasks/tasks/)

## [Timesheets API](https://help.getharvest.com/api-v2/timesheets-api)
* [Time Entries](https://help.getharvest.com/api-v2/timesheets-api/timesheets/time-entries/)

## [Projects API](https://help.getharvest.com/api-v2/projects-api)
* [Project User Assignments](https://help.getharvest.com/api-v2/projects-api/projects/user-assignments/)
* [Project Task Assignments](https://help.getharvest.com/api-v2/projects-api/projects/task-assignments/)
* [Projects](https://help.getharvest.com/api-v2/projects-api/projects/projects/)

## [Roles API](https://help.getharvest.com/api-v2/roles-api)
* [Roles](https://help.getharvest.com/api-v2/roles-api/roles/roles/)

## [Users API](https://help.getharvest.com/api-v2/users-api)
* [User Project Assignments](https://help.getharvest.com/api-v2/users-api/users/project-assignments/)
* [Users](https://help.getharvest.com/api-v2/users-api/users/users/)

## Examples
### Create client
```
ctx := context.Background()
ts := oauth2.StaticTokenSource(
	&oauth2.Token{
        AccessToken: os.Getenv("HARVEST_ACCESS_TOKEN"),
    },
)
tc := oauth2.NewClient(ctx, ts)

```

### Create service
```
service := harvest.NewHarvestClient(tc)
service.AccountId = os.Getenv("HARVEST_ACCOUNT_ID")
```

### Get organisation
```
c, _, err := service.Company.Get(ctx)
if err != nil {
    log.Error(err)
    return
}

fmt.Println("Company info")
fmt.Println(c.String())
```

### Get clients
```
clientList, _, err := service.Client.List(ctx, &harvest.ClientListOptions{})
if err != nil {
    log.Error(err)
    return
}

fmt.Println("Client list")
fmt.Println(clientList.String())
```

### Get contacts
```
contactList, _, err := service.Client.ListContacts(ctx, &harvest.ClientContactListOptions{})
if err != nil {
    log.Error(err)
    return
}

fmt.Println("Contact list")
fmt.Println(contactList.String())
```

### Get projects
```
projectList, _, err := service.Project.List(ctx, &harvest.ProjectListOptions{})
if err != nil {
    log.Error(err)
    return
}

fmt.Println("Project list")
fmt.Println(projectList.String())
```

### Get tasks
```
taskList, _, err := service.Task.List(ctx, &harvest.TaskListOptions{})
if err != nil {
    log.Error(err)
    return
}

fmt.Println("Task list")
fmt.Println(taskList.String())
```

### Get users
```
userList, _, err := service.User.List(ctx, &harvest.UserListOptions{})
if err != nil {
    log.Error(err)
    return
}

fmt.Println("User list")
fmt.Println(userList.String())
```

### Get estimates
```
estimateList, _, err := service.Estimate.List(ctx, &harvest.EstimateListOptions{})
if err != nil {
    log.Error(err)
    return
}

fmt.Println("Estimate list")
fmt.Println(estimateList.String())
```

### Get invoices
```
invoiceList, _, err := service.Invoice.List(ctx, &harvest.InvoiceListOptions{})
if err != nil {
    log.Error(err)
    return
}

fmt.Println("Invoice list")
fmt.Println(invoiceList.String())
```

### Get roles
```
roleList, _, err := service.Role.List(ctx, &harvest.RoleListOptions{})
if err != nil {
    log.Error(err)
    return
}

fmt.Println("Role list")
fmt.Println(roleList.String())
```

