go-aqbanking
============

[![Build Status](https://travis-ci.org/umsatz/go-aqbanking.svg)](https://travis-ci.org/umsatz/go-aqbanking)

golang wrapper for aqbanking C library

go run !(*_test).go

## A note about parallalization

- You may not have more than one instance of AQBanking per process.
  I'm not sure yet if this is a limit of aqbanking, or rather my fault.
  But Aqbanking panics here if you try: src/libs/aqbanking/banking_init.c:184

- You may not have more than one instance of AQBanking running with the same directory.
  E.g. 2 processes with one AQBanking instance created via DefaultAQBanking()
  will break as well.
  This seems to be due to the fact that aqbanking extensivly writes files during its operation

## TODO

- Test for memory leaks
- Write-Support
  - Transactions
    - Standing Orders
    - Foreign Transfers
    - Investment Transfers