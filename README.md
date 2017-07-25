# go-holiday-calculations
This is a test bed Go Language program for calculating holidays in various countries, loading them into a struct() and then marshalling out the data in: 
  - json
  - bson
  - yaml
  - xml
# Intent
To be a learning platform for myself to explore advanced Go functionality while updating some old Java code I wrote back in late 1997, early 1998 for Java World ( http://www.javaworld.com/article/2077543/learn-java/java-tip-44--calculating-holidays-and-their-observances.html )
# Features
  - struct() creation and population
  - multithreading of calculations using waitGroups 
  - marshalling data from the struct into files
# To do
  - break code up into packages
  - Finish some country holiday claculations
  - Beef up error processing
  - Allow input to change the default input directory from "X:\\go\\output"
  - Re-Architect the code since Go does not support constructor methods.