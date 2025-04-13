from os import getenv

from agno.docker.app.fastapi import FastApi
from agno.docker.app.postgres import PgVectorDb
from agno.docker.resource.image import DockerImage
from agno.docker.resources import DockerResources

from workspace.settings import ws_settings

#
# -*- Resources for the Development Environment
#

# -*- Dev image
dev_image = DockerImage(
    name=f"{ws_settings.image_repo}/{ws_settings.image_name}",
    tag=ws_settings.dev_env,
    enabled=ws_settings.build_images,
    path=str(ws_settings.ws_root),
    # Do not push images after building
    push_image=ws_settings.push_images,
)

# -*- Dev database running on port 5433:5433
dev_db = PgVectorDb(
    name=f"{ws_settings.ws_name}-db",
    pg_user="ai",
    pg_password="ai",
    pg_database="ai",
    # Connect to this db on port 5433
    host_port=5433,
)

# -*- Container environment
container_env = {
    "RUNTIME_ENV": "dev",
    # Get the OpenAI API key and Exa API key from the local environment
    "OPENAI_API_KEY": getenv("OPENAI_API_KEY"),
    "EXA_API_KEY": getenv("EXA_API_KEY"),
    # Enable monitoring
    "AGNO_MONITOR": "True",
    "AGNO_API_KEY": getenv("AGNO_API_KEY"),
    # Database configuration
    "DB_HOST": dev_db.get_db_host(),
    "DB_PORT": dev_db.get_db_port(),
    "DB_USER": dev_db.get_db_user(),
    "DB_PASS": dev_db.get_db_password(),
    "DB_DATABASE": dev_db.get_db_database(),
    # Wait for database to be available before starting the application
    "WAIT_FOR_DB": dev_db.enabled,
    # Migrate database on startup using alembic
    "MIGRATE_DB": dev_db.enabled,
}

# -*- FastApi running on port 8000:8000
dev_fastapi = FastApi(
    name=f"{ws_settings.ws_name}-api",
    image=dev_image,
    command="uvicorn api.main:app --reload",
    port_number=8000,
    debug_mode=True,
    mount_workspace=True,
    env_vars=container_env,
    use_cache=True,
    # Read secrets from secrets/dev_api_secrets.yml
    secrets_file=ws_settings.ws_root.joinpath("workspace/secrets/dev_api_secrets.yml"),
    depends_on=[dev_db],
)

# -*- Dev DockerResources
dev_docker_resources = DockerResources(
    env=ws_settings.dev_env,
    network=ws_settings.ws_name,
    apps=[dev_db, dev_fastapi],
)
