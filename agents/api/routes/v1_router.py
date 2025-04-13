from fastapi import APIRouter

from api.routes.agents import agents_router
from api.routes.playground import playground_router
from api.routes.status import status_router

v1_router = APIRouter(prefix="/v1")
v1_router.include_router(status_router)
v1_router.include_router(agents_router)
v1_router.include_router(playground_router)
