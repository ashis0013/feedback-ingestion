package repository

var createTables = []string{
	createRecordTable,
	createTenantTable,
	createSourceTable,
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

const createSourceTable = `
create table if not exists source_records
    (
        source_id varchar(255) primary key,
        source_name varchar(255),
        source_meta varchar(1024)
    )`

const createTenantTable = `
create table if not exists tenant_records
    (
        tenant_id varchar(255) primary key,
        tenant_name varchar(255),
        product_tags varchar(1024)
    )`

const insertRecord = `
    insert into feedback_records
    (record_id,source_id,tenant_id,person_name,person_email,feedback_type,feedback_content,record_lang,created_on,additional_data)
    values (:record_id,:source_id,:tenant_id,:person_name,:person_email,:feedback_type,:feedback_content,:record_lang,:created_on,:additional_data)
`

const insertTenant = `
    insert into tenant_records
    (tenant_id, tenant_name, product_tags)
    values (:tenant_id,:tenant_name,:product_tags)
`

const insertSource = `
    insert into source_records
    (source_id, source_name, source_meta)
    values (:source_id,:source_name,:source_meta)
`

const selectTenants = "select * from public.tenant_records"
const selectSources = "select * from public.source_records"
