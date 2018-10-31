require 'csv'

csv = CSV.read("./dictionary_of_synonym.tsv", col_sep: "\t", headers: false)
csv.each do  |row|
  puts row[0]
end
