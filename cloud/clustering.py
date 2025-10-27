import numpy as np
from sklearn.metrics.pairwise import cosine_distances
import uuid

class FaceClustering:
	def __init__(self, distance_threshold=0.6):
		self.distance_threshold = distance_threshold
	
	def find_cluster(self, embedding, cluster_embeddings):
		if not cluster_embeddings:
			return None
		
		emb_array = np.array(embedding).reshape(1, -1)
		cluster_array = np.array(cluster_embeddings).reshape(len(cluster_embeddings), -1)
		
		distances = cosine_distances(emb_array, cluster_array)[0]
		min_distance = np.min(distances)
		
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
