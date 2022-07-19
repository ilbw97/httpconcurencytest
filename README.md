# httptest
HTTP TEST TOOL USING GOROUTINE.
YOU CAN MAKE HTTP REQUEST CONCURENTLY.

## HOW TO USE?
1. run test <br>
  ```./httptest -protocol http -host wordpress.jam10000bo.com -method post -port 8099 -path /cloud2team -count 102 <br>```
  ==> you can make '102' request to http://wordpress.jam10000bo.com:8099/cloud2team at one time!! <br>
  YOU CAN MAKE HTTP REQUEST WITH OPTIONS.
2. clean up logs <br>
  just run cleanup.sh
  ```./cleanup.sh```

## YOU CAN GIVE OPTIONS BELOW.
### PROTOCOL
http or https<br>
EX) <u>**https**</u>://github.com/ilbw97/httptest/new/master?readme=1

### HOST
EX) https://<u>**github.com**</u>/ilbw97/httptest/new/master?readme=1

### METHOD
USING ONLY GET / PUT / POST / UPDATE

### PATH
URL PATH FOR REQUEST. 
EX) https://github.com <u>**/ilbw97/httptest/new/master?readme=1**</u>

### PORT
PLASE ENTER ONLY NUMERIC

### COUNT
count for goroutines you want to make.
PLASE ENTER ONLY NUMERIC


