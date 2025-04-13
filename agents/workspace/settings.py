from pathlib import Path

from agno.workspace.settings import WorkspaceSettings

#
# We define workspace settings using a WorkspaceSettings object
# these values can also be set using environment variables
# Import them into your project using `from workspace.settings import ws_settings`
#
ws_settings = WorkspaceSettings(
    # Workspace name
    ws_name="agents",
    # Path to the workspace root
    ws_root=Path(__file__).parent.parent.resolve(),
    # -*- Workspace Environments
    dev_env="dev",
    prd_env="prd",
    # -*- Workspace Keys
    # default env for `agno ws` commands
    default_env="dev",
    # -*- Image Settings
    # Repository for images
    image_repo="agnohq",
    # 'Name:tag' for the image
    image_name="agent-api",
    # Build images locally
    build_images=False,
    # Push images after building
    push_images=False,
    # Skip cache when building images
    skip_image_cache=False,
    # Force pull images in FROM
    force_pull_images=False,
    # -*- AWS settings
    # Region for AWS resources
    aws_region="us-east-1",
    # Availability Zones for AWS resources
    aws_az1="us-east-1a",
    aws_az2="us-east-1b",
    # Subnets for AWS resources
    # aws_subnet_ids=["subnet-xyz", "subnet-xyz"],
    # Security Groups for AWS resources
    # aws_security_group_ids=["sg-xyz", "sg-xyz"],
)
