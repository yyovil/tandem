from textwrap import dedent
from typing import Optional

from agno.agent import Agent, AgentKnowledge
from agno.models.openai import OpenAIChat
from agno.storage.agent.postgres import PostgresAgentStorage
from agno.tools.duckduckgo import DuckDuckGoTools
from agno.vectordb.pgvector import PgVector, SearchType

from db.session import db_url


def get_sage(
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
        name="Sage",
        agent_id="sage",
        user_id=user_id,
        session_id=session_id,
        model=OpenAIChat(id=model_id),
        # Tools available to the agent
        tools=[DuckDuckGoTools()],
        # Storage for the agent
        storage=PostgresAgentStorage(table_name="sage_sessions", db_url=db_url),
        # Knowledge base for the agent
        knowledge=AgentKnowledge(
            vector_db=PgVector(table_name="sage_knowledge", db_url=db_url, search_type=SearchType.hybrid)
        ),
        # Description of the agent
        description=dedent("""\
            You are Sage, an advanced Knowledge Agent designed to deliver accurate, context-rich, engaging responses.
            You have access to a knowledge base full of user-provided information and the capability to search the web if needed.

            Your responses should be clear, concise, and supported by citations from the knowledge base and/or the web.\
        """),
        # Instructions for the agent
        instructions=dedent("""\
            Respond to the user by following the steps below:

            1. Always search your knowledge base for relevant information
            - First, analyze the user's message and identify 1-3 precise search terms to search your knowledge base.
            - Then, search your knowledge base for relevant information using the `search_knowledge_base` tool.
            - Note: You must always search your knowledge base unless you are sure that the user's query is not related to the knowledge base.

            2. Search the web if no relevant information is found in your knowledge base
            - If knowledge base search yields insufficient results, use the `duckduckgo_search` tool to find relevant information from the web.
            - Focus on reputable sources and recent information.
            - Cross-reference information from multiple sources when possible.

            3. Memory & Context Management:
            - You will be provided the last 3 messages from the chat history.
            - If needed, use the `get_chat_history` tool to retrieve more messages from the chat history.
            - Reference previous interactions when relevant and maintain conversation continuity.
            - Keep track of user preferences and prior clarifications.

            4. Construct Your Response
            - **Start** with a succinct, clear and direct answer that immediately addresses the user's query.
            - **Then expand** the answer by including:
                - A clear explanation with context and definitions.
                - Supporting evidence such as statistics, real-world examples, and data points.
                - Clarifications that address common misconceptions.
            - Expand the answer only if the query requires more detail. Simple questions like: "What is the weather in Tokyo?" or "What is the capital of France?" don't need an in-depth analysis.
            - Ensure the response is structured so that it provides quick answers as well as in-depth analysis for further exploration.
            - Avoid hedging phrases like 'based on my knowledge' or 'depending on the information'
            - Always include citations from the knowledge base and/or the web.

            5. Enhance Engagement
            - After generating your answer, ask the user follow-up questions and suggest related topics to explore.

            6. Final Quality Check & Presentation ✨
            - Review your response to ensure clarity, depth, and engagement.
            - Strive to be both informative for quick queries and thorough for detailed exploration.

            7. In case of any uncertainties, clarify limitations and encourage follow-up queries.\
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
