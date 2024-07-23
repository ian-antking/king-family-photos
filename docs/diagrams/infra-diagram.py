from diagrams import Diagram, Cluster, Edge
from diagrams.generic.storage import Storage
from diagrams.generic.device import Tablet
from diagrams.aws.storage import SimpleStorageServiceS3 as S3
from diagrams.aws.compute import Lambda

with Diagram("King Family Photos", graph_attr={ "concentrate": "true" }):
  nas = Storage("NAS")
  photo_bucket = S3("Photo Backup")
  photo_frame = Tablet("Photo Frame")
    
  nas >> Edge(taillabel="s3 sync", minlen="2") >> photo_bucket
  
  with Cluster("Serverless App"):
    photo_display_bucket = S3("Photo Display")
    resize_photo = Lambda("resizePhoto")
    remove_photo = Lambda("removePhoto")
    
    (
      photo_bucket >> 
      Edge(headlabel="s3:ObjectCreated", style="dashed", minlen="3") >> 
      resize_photo >> Edge(taillabel="putObject", minlen="3", labelloc="t") >> 
      photo_display_bucket
    )
    
    (
      photo_bucket >> 
      Edge(headlabel="s3:ObjectRemoved", style="dashed", minlen="3") >> 
      remove_photo >> Edge(taillabel="deleteObject", minlen="3", labelloc="t") >>  
      photo_display_bucket
    )
    
    photo_display_bucket << Edge(headlabel="s3 sync", minlen="2") << photo_frame
    
  