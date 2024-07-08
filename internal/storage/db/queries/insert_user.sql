insert into users(surname, name, patronymic, address, passport)
values ($1, $2, $3, $4, $5)
returning id
