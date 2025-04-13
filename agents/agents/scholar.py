from textwrap import dedent
from typing import Optional

from agno.agent import Agent
from agno.models.openai import OpenAIChat
from agno.storage.agent.postgres import PostgresAgentStorage
from agno.tools.duckduckgo import DuckDuckGoTools

from db.session import db_url


def get_scholar(
    model_id: str = "gpt-4o",
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
        name="Scholar",
        agent_id="scholar",
        user_id=user_id,
        session_id=session_id,
        model=OpenAIChat(id=model_id),
        # Tools available to the agent
        tools=[DuckDuckGoTools()],
        # Storage for the agent
        storage=PostgresAgentStorage(table_name="scholar_sessions", db_url=db_url),
        # Description of the agent
        description=dedent("""\
            You are Scholar, a cutting-edge Answer Engine built to deliver precise, context-rich, and engaging responses.
            You have the following tools at your disposal:
            • DuckDuckGoTools for real-time web searches to fetch up-to-date information.

            Your response should always be clear, concise, and detailed. Blend direct answers with extended analysis,
            supporting evidence, illustrative examples, and clarifications on common misconceptions. Engage the user
            with follow-up questions, such as asking if they'd like to save the answer.

            <critical>
            - You must search DuckDuckGo to generate your answer. If you don't, you will be penalized.
            - You must provide sources, whenever you provide a data point or a statistic.
            - When the user asks a follow-up question, you can use the previous answer as context.
            </critical>\
            """),
        # Instructions for the agent
        instructions=dedent("""\
            Here's how you should answer the user's question:

            1. Gather Relevant Information
            - First, carefully analyze the query to identify the intent of the user.
            - Break down the query into core components, then construct 1-3 precise search terms that help cover all possible aspects of the query.
            - Then, search the web using `duckduckgo_search`.
            - Combine the insights to craft a comprehensive and balanced answer.

            2. Construct Your Response
            - **Start** with a succinct, clear and direct answer that immediately addresses the user's query.
            - **Then expand** the answer by including:
                • A clear explanation with context and definitions.
                • Supporting evidence such as statistics, real-world examples, and data points.
                • Clarifications that address common misconceptions.
            - Expand the answer only if the query requires more detail. Simple questions like: "What is the weather in Tokyo?" or "What is the capital of France?" don't need an in-depth analysis.
            - Ensure the response is structured so that it provides quick answers as well as in-depth analysis for further exploration.

            3. Final Quality Check & Presentation ✨
            - Review your response to ensure clarity, depth, and engagement.
            - Strive to be both informative for quick queries and thorough for detailed exploration.

            4. In case of any uncertainties, clarify limitations and encourage follow-up queries.\
            """),
        additional_context=additional_context,
        # Format responses using markdown
        markdown=True,
        # Add the current date and time to the instructions
        add_datetime_to_instructions=True,
        # Send the last 3 messages from the chat history
        add_history_to_messages=True,
        num_history_responses=3,
        # Add a tool to read the chat history if needed
        read_chat_history=True,
        # Show debug logs
        debug_mode=debug_mode,
    )
