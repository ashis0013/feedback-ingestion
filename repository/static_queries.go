package repository

var createTables = []string{
	createRecordTable,
	createTenantTable,
}

const createRecordTable = `
create table if not exists feedback_records
    (
        record_id varchar(255) primary key,
        source_id varchar(255),
        tenant_id varchar(255),
        person_name varchar(255),
        person_email varchar(255),
        feedback_type varchar(50),
        feedback_content varchar(4095),
        record_lang varchar(50),
        created_on timestamp,
        additional_data varchar(32768)
    )`

const createTenantTable = `
create table if not exists tenant_records
    (
        tenant_id integer primary key,
        tenant_name varchar(255)
    )`

const insertRecord = `
    insert into feedback_records
    (record_id,source_id,tenant_id,person_name,person_email,feedback_type,feedback_content,record_lang,created_on,additional_data)
    values (:record_id,:source_id,:tenant_id,:person_name,:person_email,:feedback_type,:feedback_content,:record_lang,:created_on,:additional_data)
`
