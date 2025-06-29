from typing import AsyncGenerator, List, Optional, AsyncIterator
from utils.models import Model
from agno.run.response import RunResponse
from agno.agent import Agent
from agno.media import File
from fastapi import APIRouter, HTTPException, status
from fastapi.responses import StreamingResponse
from pydantic import BaseModel

from agents.operator import AgentType, get_agent, get_available_agents
from utils.log import logger

######################################################
## Router for the Agent Interface
######################################################

agents_router = APIRouter(prefix="/agents", tags=["Agents"])


@agents_router.get("", response_model=List[str])
async def list_agents():
    """
    Returns a list of all available agent IDs.

    Returns:
        List[str]: List of agent identifiers
    """
    return get_available_agents()


async def chat_response_streamer(agent: Agent, message: str, attachments: Optional[List[File]]) -> AsyncGenerator:
    """
    Stream agent responses chunk by chunk.

    Args:
        agent: The agent instance to interact with
        message: User message to process
        attachments: User attachments to process
        attachments: User attachments to process
    Yields:
        chunks serialised in JSON, from the agent response.
        chunks serialised in JSON, from the agent response.
    """

    run_responses: AsyncIterator[RunResponse]

    if (attachments is not None) and len(attachments) > 0:
        attachments = list(
            map(
                lambda attachment: File(
                    url=attachment.url or None,
                    content=attachment.content or None,
                    mime_type=attachment.mime_type or "text/plain",
                    filepath=attachment.filepath or None,
                ),
                attachments,
            )
        )
        run_responses = await agent.arun(message, stream=True, files=attachments)
    else:
        run_responses = await agent.arun(message, stream=True)

    async for chunk in run_responses:
        # messages are separated by a pair of newline characters.
        yield chunk.to_json() + "\n\n"


# extend this schema to support the local files and URLs.
class RunRequest(BaseModel):
    """Request model for an running an agent"""

    message: str
    stream: Optional[bool] = True
    model_id: Model = Model.GEMINI_2_5_FLASH_PREVIEW_04_17.value
    user_id: Optional[str] = None
    session_id: Optional[str] = None
    attachments: Optional[List[File]] = None


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
            model_id=body.model_id,
            agent_id=agent_id,
            user_id=body.user_id,
            session_id=body.session_id,
        )
    except Exception as e:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail=f"Agent not found: {str(e)}")

    if body.stream:
        return StreamingResponse(
            chat_response_streamer(agent, body.message, body.attachments),
            media_type="text/event-stream",
        )
    else:
        response = await agent.arun(body.message, stream=False, files=body.attachments)
        # response.content only contains the text response from the Agent.
        # For advanced use cases, we should yield the entire response
        # that contains the tool calls and intermediate steps.
        # response.content only contains the text response from the Agent.
        # For advanced use cases, we should yield the entire response
        # that contains the tool calls and intermediate steps.
        # response.content only contains the text response from the Agent.
        # For advanced use cases, we should yield the entire response
        # that contains the tool calls and intermediate steps.
        return response.content
