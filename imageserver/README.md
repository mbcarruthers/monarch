# imageserver

Dead simple image server. Put anything in assets, build the docker-compose file,
and as long as they are images they will be served in bytes. No sizing or anything
on the backend other than what the image is sent in bytes. 

Should work for anything that is an image but it is working as of now of of .jpg's
but ...

Todo: Test Image formats and look into setting the w.Header() for images.

Look into seeing if it would be possible to only accept certain formats. May in 
the admin UI or something. 

See if this type of thing would be compatible with Lura.