echo Compressing distribution 
zip -r dist.zip build/ scripts/ Makefile
echo Pushing dist to server
gcloud beta compute scp -q --zone "us-central1-a" --project "spatiumsocialis" ./dist.zip spatium-prod:~/
echo Pulling service images and starting services containers
gcloud beta compute ssh  'spatium-prod' --zone "us-central1-a"  --project "spatiumsocialis"  --command 'unzip -o ./dist.zip; make pull; make start env=prod; exit'
echo Removing dist.zip
rm ./dist.zip
