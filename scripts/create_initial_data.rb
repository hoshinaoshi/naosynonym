require 'csv'
require 'aws-sdk-dynamodb'

ddb = Aws::DynamoDB::Client.new(
  endpoint: 'http://localhost:4569',
  region: 'us-east-1'
)


csv = CSV.read("./dictionary_of_synonym.tsv", col_sep: "\t", headers: false)
csv.each do  |row|
  puts row[0]
  puts row[1]
  puts row[2]
  puts row[3]
  puts row[4]
  puts row[5]

  ddb.put_item(
    table_name: "test",
    item:  {
      tag:             row[0],
      readings:        row[1],
      readings2:       row[2],
      part_of_speech:  row[3],
      conjugation:     row[4],
      synonyms:        row[5]
    }
  )
end
