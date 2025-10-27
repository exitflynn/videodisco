import os
import numpy as np
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from sqlalchemy import create_engine, func
from sqlalchemy.orm import sessionmaker
from models import Base, FaceGroup, Embedding
from clustering import FaceClustering
import uuid
import struct

app = FastAPI(title="VideoDisco Cloud Clustering")

DATABASE_URL = os.getenv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/videodisco_cloud")
engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(bind=engine)

Base.metadata.create_all(bind=engine)
clusterer = FaceClustering(distance_threshold=0.6)

class ClusterRequest(BaseModel):
	image_id: str
	embedding: list
	metadata: dict = {}

class ClusterResponse(BaseModel):
	image_id: str
	group_id: str
	is_new_group: bool

@app.post("/cluster/add", response_model=ClusterResponse)
def add_to_cluster(request: ClusterRequest):
	db = SessionLocal()
	try:
		emb_array = np.array(request.embedding, dtype=np.float32)
		
		groups = db.query(FaceGroup).all()
		assigned_group = None
		min_dist = float('inf')
		
		for group in groups:
			group_embeddings = [np.frombuffer(e.embedding, dtype=np.float32) for e in group.embeddings]
			if not group_embeddings:
				continue
			
			emb_test = emb_array.reshape(1, -1)
			group_arr = np.array(group_embeddings)
			distances = cosine_distances(emb_test, group_arr)[0]
			min_d = np.min(distances)
			
			if min_d < clusterer.distance_threshold and min_d < min_dist:
				min_dist = min_d
				assigned_group = group
		
		if assigned_group:
			group_id = assigned_group.group_id
			is_new = False
		else:
			group_id = str(uuid.uuid4())
			new_group = FaceGroup(group_id=group_id)
			db.add(new_group)
			db.flush()
			is_new = True
		
		embedding_bytes = emb_array.tobytes()
		embedding = Embedding(
			image_id=request.image_id,
			group_id=group_id,
			embedding=embedding_bytes
		)
		db.add(embedding)
		db.commit()
		
		return ClusterResponse(
			image_id=request.image_id,
			group_id=group_id,
			is_new_group=is_new
		)
	except Exception as e:
		db.rollback()
		raise HTTPException(status_code=400, detail=str(e))
	finally:
		db.close()

@app.get("/health")
def health():
	return {"status": "ok"}

@app.get("/clusters")
def list_clusters():
	db = SessionLocal()
	try:
		groups = db.query(FaceGroup).all()
		return {
			"total_groups": len(groups),
			"groups": [
				{
					"group_id": g.group_id,
					"count": len(g.embeddings),
					"created_at": g.created_at.isoformat()
				}
				for g in groups
			]
		}
	finally:
		db.close()

if __name__ == "__main__":
	import uvicorn
	uvicorn.run(app, host="127.0.0.1", port=5002)

from sklearn.metrics.pairwise import cosine_distances
