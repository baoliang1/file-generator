# file-generator

docker build -t go-file-generator .
docker run -p 8080:8080 go-file-generator

kubectl apply -f file-generator.yaml
kubectl get service file-generator-service

