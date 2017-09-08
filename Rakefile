namespace :app do 
	
end

namespace :db do
  task :migrate do
    sh "rm -rf data && mkdir data && chmod -R 775 data/"
    sh "migrate -url sqlite3://data/locsvc.db -path ./migrations/ up" # resetup db
    sh "chmod 775 data/locsvc.db"
  end
end