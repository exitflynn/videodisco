import numpy as np
import uuid

class FaceClustering:
	def __init__(self, distance_threshold=0.6):
		self.distance_threshold = distance_threshold
	
	def cosine_distance(self, emb1, emb2):
		emb1 = np.array(emb1, dtype=np.float32)
		emb2 = np.array(emb2, dtype=np.float32)
		
		dot_product = np.dot(emb1, emb2)
		norm1 = np.linalg.norm(emb1)
		norm2 = np.linalg.norm(emb2)
		
		if norm1 == 0 or norm2 == 0:
			return 1.0
		
		cosine_sim = dot_product / (norm1 * norm2)
		return 1.0 - cosine_sim
	
	def find_cluster(self, embedding, cluster_embeddings):
		if not cluster_embeddings:
			return None
		
		min_distance = float('inf')
		for cluster_emb in cluster_embeddings:
			dist = self.cosine_distance(embedding, cluster_emb)
			if dist < min_distance:
				min_distance = dist
		
		if min_distance < self.distance_threshold:
			return min_distance
		return None
	
	def compute_cluster_center(self, embeddings):
		if not embeddings:
			return None
		return np.mean(embeddings, axis=0)
	
	def assign_to_cluster(self, embedding):
		dist = self.find_cluster(embedding, [embedding])
		if dist is None:
			return str(uuid.uuid4())
		return None
