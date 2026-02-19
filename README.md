# SQL Test suit for SQL adapters (wip)

swiss army lib: sql-adapter-test: `sqltest`

`sqltest.SetupTest`:

instantiates a test-container in mysql, with proper cleanup:

- instantiates NexusTest
- instantiates CpoolTest (using the returned nexus instance)
- instantiates SQLTest

returns:

- mysql connection string
- cpool instance

`sqltest.SetupNexusTest`:

uses user input to create and populate nexus db

returns nexus connection string and start a nexus instance using swiss army (mocking all RSA operations)

`sqltest.SetupCPoolTest`:

accpets a nexus instance to provide a connection pool

`sqltest.SetupSQLTest`:

accepts:

- context (for timeouts and cancelation)
- testing.T (for cleanup configuration)
- migration version (0 is the latest one)
- populateFunction: `func() sqltest.PopulateData`  (`sqltest.PopulateFunc`). might be one of the following:
    - `sqltest.PopulateDefault` →populate all tables with default data for the selected migration (we need to provide a `data-soure:migration-verion` somewhere so that developers know which data to expect). save developing time, but increase testing processing time
    - `sqltest.PopulateEmpty` → dont populate anything, (useful for insertion-only operations, might save us some time)
    - custom populate function:
        
        ```go
        // pre defined in sqltest (only for reference) ---- inside sqltest pkg
        type PopulateData struct {
        	Rows []string   //  {id name description price}
        	Data [][]string // {{1  product desc1    10},...}
        }
        
        func NewPopulate(tableName ...string) PopulateData
        func (*PopulateData) Insert(values ...string)*PopulateData
        // end of predefined -------------------------------------------------
        
        // actual usage for developers:
        func() PopulateData {
        	return NewPopulate("id", "name", "decription", "price")
        		.Insert("1", "product1", "description1", "1")
        		.Insert("2", "product2", "description2", "10")
        }
        ```
        

returns:

- mysql connection string
- error:
    - container initialization
    - migration error
    - populate error
    - nil

it’s up to the developer connect to that instance (treat it like a normal mysql connection)