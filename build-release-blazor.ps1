cd BlazorApp1
dotnet publish -c Release -o ./publish
cd ..
cd .\BlazorApp2
dotnet publish -c Release -o ./publish
cd ..
 
go build ./cmd/httpserver/
