from enum import Enum
from typing import List, Optional
from agents.reconnoiter import get_reconnoiter
from utils.models import Model


class AgentType(Enum):
    RECONNOITER = "Mr. Burnham"


def get_available_agents() -> List[str]:
    """Returns a list of all available agent IDs."""
    return [agent.value for agent in AgentType]


def get_agent(
    model_id: str = Model.GEMINI_2_5_FLASH_PREVIEW_04_17.value,
    agent_id: Optional[AgentType] = None,
    user_id: Optional[str] = None,
    session_id: Optional[str] = None,
    debug_mode: bool = True,
):
    match agent_id:
        case AgentType.RECONNOITER:
            return get_reconnoiter(
                model_id=model_id,
                user_id=user_id,
                session_id=session_id,
                debug_mode=debug_mode,
            )
        case _:
            raise Exception(f"Agent {agent_id} is not available.")
