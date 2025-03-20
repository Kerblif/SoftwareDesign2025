define setup_env
    $(eval include .env)
    $(eval export)
endef

generate-db:
	$(call setup_env)
	cd ./HSEFinance.Lib && dotnet ef migrations add InitialCreate && dotnet ef database update

run:
	$(call setup_env)
	cd ./HSEFinance.ConsoleApp && dotnet run
