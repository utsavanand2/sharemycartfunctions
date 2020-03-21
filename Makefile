userCreate:
	cd authEvent && \
gcloud functions deploy UserCreated \
--trigger-event providers/firebase.auth/eventTypes/user.create \
--trigger-resource collabshop19 \
--runtime go113

updateNeed:
	cd updateItemInFriendFromNeed && \
gcloud functions deploy UpdateListToAddNeed \
--trigger-http \
--allow-unauthenticated \
--runtime go113