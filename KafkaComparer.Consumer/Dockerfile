#BUILD
FROM microsoft/dotnet:2.1.302-sdk AS build
WORKDIR /src

COPY . .
RUN dotnet restore -nowarn:msb3202,nu1503
RUN dotnet build -c Release -o /app

FROM build AS publish
RUN dotnet publish -c Release -o /app

FROM microsoft/dotnet:2.1-aspnetcore-runtime AS runtime
WORKDIR /app
COPY --from=publish /app .
ENTRYPOINT ["dotnet", "KafkaComparer.Consumer.dll"]