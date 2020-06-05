# Arguments
# $1: Google GCR hostname
# $2: Google Project ID
# $3: Build package directory
# $4: Docker build context

source ./scripts/dockerhost.sh
docker build -t $1/$2/deps:latest -f $3/deps.Dockerfile $4
