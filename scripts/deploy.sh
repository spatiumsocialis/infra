gcloud beta compute scp --zone "us-central1-a" --project "spatiumsocialis" ./docker-compose.yml ./docker-compose.prod.yml Makefile spatium-prod:~/
gcloud beta compute ssh  'spatium-prod' --zone "us-central1-a"  --project "spatiumsocialis"  --command 'make pull; make start env=prod; exit'

