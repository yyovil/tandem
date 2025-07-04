from typing import Optional
from agno.agent import Agent
from agno.models.google import Gemini
from agno.tools.docker import DockerTools
from textwrap import dedent
from utils.models import Model


def get_reconnoiter(
    model_id: str = Model.GEMINI_2_5_FLASH_PREVIEW_04_17.value,
    user_id: Optional[str] = None,
    session_id: Optional[str] = None,
    debug_mode: bool = True,
) -> Agent:
    additional_context = ""
    if user_id:
        additional_context += "<context>"
        additional_context += f"You are interacting with the user: {user_id}"
        additional_context += "</context>"

    return Agent(
        name="Frederick Russell Burnham",
        agent_id="Mr. Burnham",
        user_id=user_id,
        session_id=session_id,
        model=Gemini(id=model_id),
        description=dedent(
            """You are an OffSec PEN-300 certified penetration tester with extensive experience in reconnaissance."""
        ),
        goal=dedent(
            """Your goal is to assist the user during the reconnaissance phase of a pentest and finish the tasks assigned."""
        ),
        instructions=[
            "Be concise and clear.",
            "Use kali linux cli tools for reconnaissance.",
            "Use the kali:withtools image to spawn a new container.",
            "Connect to host docker network.",
            "get a bash shell in the container.",
            "Reuse any kali:withtools image containers if they are available.",
            "Progressively exec the bash cmds in the docker container.",
            "Always ask for clarification if certain things aren't clear to you.",
            "Always put the scanning results in a txt file with this name scheme: {tool_used}_{scan_type}.txt, using the redirection operator in bash.",
        ],
        tools=[
            DockerTools(
                enable_container_management=True,
                enable_image_management=True,
                enable_network_management=True,
                enable_volume_management=True,

            )
        ],
        show_tool_calls=True,
        markdown=True,
        add_datetime_to_instructions=True,
        add_history_to_messages=True,
        add_name_to_instructions=True,
        stream=True,
        read_chat_history=True,
        add_state_in_messages=True,
        debug_mode=debug_mode,
        additional_context=additional_context,
    )
