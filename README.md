go-aqbanking
============

[![CircleCI](https://circleci.com/gh/umsatz/go-aqbanking/tree/master.svg?style=shield)](https://circleci.com/gh/umsatz/go-aqbanking/tree/master)

golang wrapper for aqbanking C library

go run !(*_test).go

## Dependencies

If you are running Ubuntu 18.04:

    sudo add-apt-repository ppa:ingo/gnucash
    sudo apt-get update
    sudo apt-get install libgwenhywfar-core-dev libaqbanking43

## A note about parallelization

- You may not have more than one instance of AQBanking per process.
  I'm not sure yet if this is a limit of aqbanking, or rather my fault.
  Either way, Aqbanking panics here if you try: src/libs/aqbanking/banking_init.c:184

- You may not have more than one instance of AQBanking running with the same directory.
  E.g. 2 processes with one AQBanking instance created via DefaultAQBanking()
  will break as well.
  This seems to be due to the fact that aqbanking extensively writes files during its operation.

## TODO

- Test for memory leaks
- Write-Support
  - Transactions
    - Standing Orders
    - Foreign Transfers
    - Investment Transfers
