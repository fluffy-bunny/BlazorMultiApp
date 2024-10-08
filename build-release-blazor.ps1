cd BlazorApp1
dotnet publish -c Release -o ./publish
cd ..
# Delete existing wwwroot folder
Remove-Item -Path ./cmd/httpserver/static/app1/wwwroot -Recurse -Force

# Copy new wwwroot folder
Copy-Item -Path ./BlazorApp1/publish/wwwroot -Destination ./cmd/httpserver/static/app1/wwwroot -Recurse

# Rename index.html to index_template.html
Rename-Item -Path ./cmd/httpserver/static/app1/wwwroot/index.html -NewName index_template.html

# Modify base href in index_template.html
$content = Get-Content -Path ./cmd/httpserver/static/app1/wwwroot/index_template.html -Raw
$content = $content -replace '<base href="/" />', '<base href="/app1/" />'
$content | Set-Content -Path ./cmd/httpserver/static/app1/wwwroot/index_template.html


cd .\BlazorApp2
dotnet publish -c Release -o ./publish
cd ..


# Delete existing wwwroot folder
Remove-Item -Path ./cmd/httpserver/static/app2/wwwroot -Recurse -Force

# Copy new wwwroot folder
Copy-Item -Path ./BlazorApp2/publish/wwwroot -Destination ./cmd/httpserver/static/app2/wwwroot -Recurse

# Rename index.html to index_template.html
Rename-Item -Path ./cmd/httpserver/static/app2/wwwroot/index.html -NewName index_template.html

# Modify base href in index_template.html
$content = Get-Content -Path ./cmd/httpserver/static/app2/wwwroot/index_template.html -Raw
$content = $content -replace '<base href="/" />', '<base href="/app2/" />'
$content | Set-Content -Path ./cmd/httpserver/static/app2/wwwroot/index_template.html

go build ./cmd/httpserver/
