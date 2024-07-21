## About
Welcome to the project **Cron Parser**. 

Given a cron expression string, it expands each time field to show the times at which it will run.

#### Scope:
 - It runs as command line application with cron string as a single string argument:
   ```
   ~$ your-program "cron-expression"
   ```
 - Cron expression should contain **6** fields(command included) separated by space.
 
 - Supports standard **Unix** based cron expressions with **5** time fields and **4** special characters:

    ```
    Field Name      Allowed Values    Allowed Special Characters
    ----------      --------------    --------------------------
    Minutes              0-59              * / , -
    Hours                0-23              * / , -
    Day of month         1-31              * / , - 
    Month                1-12 JAN-DEC      * / , -
    Day of week          0-6, SUN-SAT      * / , - 
    ```

   
   **Special Character Usage:**
    - **Asterisk** ( * ): Matches all possible values for the field. 
        ```
        eg: * * * * *

        Every minute of every hour of every day of the month, every month, every day of the week
        ```
    - **Comma** ( , ): Specifies a list of values.
        ```
        eg: 5,15,25 * * * *
        
        At minutes 5, 15, and 25 of every hour of every day
        ```
    - **Hyphen**   ( - ): Specifies a range of values.
        ```
        eg: 0 0-5 * * *
        
        Every hour from 12 AM to 5 AM
        ```
    - **Slash**    ( / ): Specifies intervals.
        ```
        eg: */5 * * * *
        
        Every 5 minutes
        ```


## Usage:

### Run on local machine:

***Pre-requisite:*** Install go, if not, install from [here](https://go.dev/doc/install).
1. Clone and install dependencies:
    ```
    ~$ git clone https://github.com/SravanTurbo/cron-parser.git
    ~$ go mod tidy
    ```
2. Test the repositord code:
    ```
    ~$ cd <repo>
    ~$ cd pkg/cronparser
    ~$ go test
    ```
3. Run using repository code:
    ```
    ~$ cd <repo>
    ~$ go run cmd/main.go <cron-expression>
    ~$ go run cmd/main.go "*/15 0 1,15 * 1-5 /usr/bin/find"   --> example
    ```
4. Save & Run with binary:
    ```
    ~$ cd <repo>
    ~$ go build -o ./bin/cron-parser cmd/main.go
    ~$ ./bin/cron-parser <cron-expression>
    ~$ ./bin/cron-parser "*/15 0 1,15 * 1-5 /usr/bin/find"    --> example

### As a module in your project: (TODO)

1. Add pkg to your repository:

    ```
    ~$ go get github.com/sravan-turbo/cron-parser
    ```
2. Import pkg in repository file and use:

    ```
    import github.com/sravan-turbo/cron-parser/pkg/cronparser
    ```

    ```
    cronExpr := "*/15 0 1,15 * 1-5 /usr/bin/find"
    schedule, err := cronparser.Parse(cronExpr)
    ```