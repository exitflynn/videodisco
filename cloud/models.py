from sqlalchemy import create_engine, Column, Integer, String, Float, DateTime, LargeBinary, ForeignKey, Enum as SQLEnum
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import relationship
from datetime import datetime
import enum

Base = declarative_base()

class FaceGroup(Base):
	__tablename__ = "face_groups"
	
	id = Column(Integer, primary_key=True)
	group_id = Column(String, unique=True, index=True)
	created_at = Column(DateTime, default=datetime.utcnow)
	updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)
	
	embeddings = relationship("Embedding", back_populates="group", cascade="all, delete-orphan")
	
	def __repr__(self):
		return f"<FaceGroup {self.group_id}>"


class Embedding(Base):
	__tablename__ = "embeddings"
	
	id = Column(Integer, primary_key=True)
	image_id = Column(String, unique=True, index=True)
	group_id = Column(String, ForeignKey("face_groups.group_id"))
	embedding = Column(LargeBinary)
	created_at = Column(DateTime, default=datetime.utcnow)
	
	group = relationship("FaceGroup", back_populates="embeddings")
	
	def __repr__(self):
		return f"<Embedding {self.image_id}>"
