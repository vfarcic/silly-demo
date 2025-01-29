import common

SecretName = common.Name if common.Values.db.provider != "cnpg" else "${common.Name}-app"
Name = "app" if common.Values.db.provider == "cnpg" else common.Name

CNPG = {
    apiVersion = "postgresql.cnpg.io/v1"
    kind = "Cluster"
    metadata = common.Metadata
    spec = {
        instances = 1
        storage.size = "1Gi"
    }
}

Schema = {
    apiVersion = "db.atlasgo.io/v1alpha1"
    kind = "AtlasSchema"
    metadata = common.Metadata
    spec = {
        urlFrom.secretKeyRef = { name = SecretName, key = "uri" }
        schema.sql = """\
create table videos (
    id varchar(50) not null,
    title text,
    primary key (id)
);
create table comments (
    id serial,
    video_id varchar(50) not null,
    description text not null,
    primary key (id),
    CONSTRAINT fk_videos FOREIGN KEY(video_id) REFERENCES videos(id)
);
"""
    }
}