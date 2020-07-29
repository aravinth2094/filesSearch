# Files Search
### A program that searches for files in a given directory tree with a required modified time

## Ensure
* The [search property file](search.properties) is in the same directory as `filesSearch` binary
* The property file contains key-value pairs in each line
    * Key = The directory tree you want to start the search in
        * Can be relative path or full path
    * Value = The last modified date you want filter by
        * Can be either of `today` or `YYYY-MM-DD` formatted date
    * See [Example](search.properties)

## Remember
* Output will be written in the same directory from where this software is run
    * Example: output_2020-07-29_18-10-58.csv
* The output contains the following fields
    * FILE_NAME - File name with full path
    * FILE_MODIFIED_TIME - Last modified time of the file in local timezone
    * FILE_SIZE - File Size in bytes
    * DIRECTORY_FILES_COUNT - Total number of files in the directory of the file mentioned in the current record

## Sample Console Log
```
aravinth@Aravinths-MacBook-Pro FilesSearchGo % ./dist/filesSearch_darwin-amd64.out 
2020/07/29 18:10:58 ./dist/filesSearch_darwin-amd64.out Version 1.0.0
2020/07/29 18:10:58 Loading configuration from search.properties
2020/07/29 18:10:58 Prepping output file: output_2020-07-29_18-10-58.csv
2020/07/29 18:10:58 Searching...
2020/07/29 18:10:58 Directory: /Users/aravinth/Projects => Date: 2020-07-29
2020/07/29 18:10:59 62 / 1211 found in 33 subdirectories in 15.201903ms
2020/07/29 18:10:59 Program Completed
```

## License
This software is licensed under [MIT License](LICENSE). Feel free to use it in anyway you like