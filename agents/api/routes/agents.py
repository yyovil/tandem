from enum import Enum
from typing import AsyncGenerator, List, Optional

from agno.agent import Agent
from fastapi import APIRouter, HTTPException, status
from fastapi.responses import StreamingResponse
from pydantic import BaseModel

from agents.operator import AgentType, get_agent, get_available_agents
from utils.log import logger

######################################################
## Router for the Agent Interface
######################################################

agents_router = APIRouter(prefix="/agents", tags=["Agents"])


class Model(str, Enum):
    gpt_4o = "gpt-4o"
    o3_mini = "o3-mini"


@agents_router.get("", response_model=List[str])
async def list_agents():
    """
    Returns a list of all available agent IDs.

    Returns:
        List[str]: List of agent identifiers
    """
    return get_available_agents()


async def chat_response_streamer(agent: Agent, message: str) -> AsyncGenerator:
    """
    Stream agent responses chunk by chunk.

    Args:
        agent: The agent instance to interact with
        message: User message to process

    Yields:
        Text chunks from the agent response
    """
    run_response = await agent.arun(message, stream=True)
    async for chunk in run_response:
        # chunk.content only contains the text response from the Agent.
        # For advanced use cases, we should yield the entire chunk
        # that contains the tool calls and intermediate steps.
        yield chunk.content


class RunRequest(BaseModel):
    """Request model for an running an agent"""

    message: str
    stream: bool = True
    model: Model = Model.gpt_4o
    user_id: Optional[str] = None
    session_id: Optional[str] = None


@agents_router.post("/{agent_id}/runs", status_code=status.HTTP_200_OK)
async def run_agent(agent_id: AgentType, body: RunRequest):
    """
    Sends a message to a specific agent and returns the response.

    Args:
        agent_id: The ID of the agent to interact with
        body: Request parameters including the message

    Returns:
        Either a streaming response or the complete agent response
    """
    logger.debug(f"RunRequest: {body}")

    try:
        agent: Agent = get_agent(
            model_id=body.model.value,
            agent_id=agent_id,
            user_id=body.user_id,
            session_id=body.session_id,
        )
    except Exception as e:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail=f"Agent not found: {str(e)}")

    if body.stream:
        return StreamingResponse(
            chat_response_streamer(agent, body.message),
            media_type="text/event-stream",
        )
    else:
        response = await agent.arun(body.message, stream=False)
        # response.content only contains the text response from the Agent.
        # For advanced use cases, we should yield the entire response
        # that contains the tool calls and intermediate steps.
        return response.content
