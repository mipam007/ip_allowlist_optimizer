# ip_allowlist_optimizer
take IP list and creates subnets

** create ip list with subnets
*** seems to eat excel correctly
*** sheet needs to have a name iplist
*** creating subneted IP's as follows
./ip_allowlist_optimizer alowlist.xlsx
or ./ip_allowlist_optimizer alowlist.xlsx | pbcopy 
or ./ip_allowlist_optimizer alowlist.xlsx | wc -l, etc
** build new release
*** go build -o ip_allowlist_optimizer main.go

