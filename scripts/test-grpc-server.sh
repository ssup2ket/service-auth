#!/bin/bash
# -- Variables --
EXIT_CODE=0
RANDOM_VALUE=$(shuf -i 1000-99999 -n 1)
ERROR_CODE=""
ERROR_MESSAGE=""

ADMIN_ID=""
ADMIN_ACCESS_TOKEN=""
ADMIN_REFRESH_TOKEN=""
ADMIN_LOGIN_ID="admin$RANDOM_VALUE"
ADMIN_PASSWD="admin$RANDOM_VALUE"
ADMIN_PHONE="000-0000-0000"
ADMIN_PHONE2="111-1111-1111"
ADMIN_EMAIL="admin$RANDOM_VALUE@test.com"

USER_ID=""
USER_ACCESS_TOKEN=""
USER_REFRESH_TOKEN=""
USER_LOGIN_ID="user$RANDOM_VALUE"
USER_PASSWD="user$RANDOM_VALUE"
USER_PHONE="222-2222-2222"
USER_PHONE2="333-3333-3333"
USER_EMAIL="user$RANDOM_VALUE@test.com"

# -- Test cases --
## Create users
# Admin
echo "-- Create admin user start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -d "{\"loginId\": \"$ADMIN_LOGIN_ID\", \"password\": \"$ADMIN_PASSWD\", \"phone\": \"$ADMIN_PHONE\", \"email\": \"$ADMIN_EMAIL\", \"role\" : \"admin\"}" \
  localhost:9090 User/CreateUser)
if [ $? != 0 ] ; then
  EXIT_CODE=1
fi
echo Response : $RESPONSE
echo "-- Create admin user end --"

# User
echo "-- Create user start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -d "{\"loginId\": \"$USER_LOGIN_ID\", \"password\": \"$USER_PASSWD\", \"phone\": \"$USER_PHONE\", \"email\": \"$USER_EMAIL\", \"role\" : \"user\"}" \
  localhost:9090 User/CreateUser)
if [ $? != 0 ] ; then
  EXIT_CODE=1
fi
echo Response : $RESPONSE
echo "-- Create user end --"

## Login
# Admin
echo "-- Login admin user token start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -H "username:$ADMIN_LOGIN_ID" -H "password:$ADMIN_PASSWD" \
  localhost:9090 Token/LoginToken)
if [ $? != 0 ] ; then
  EXIT_CODE=1
else
  ADMIN_ACCESS_TOKEN=$(jq -r .accessToken.token <<< "$RESPONSE")
  ADMIN_REFRESH_TOKEN=$(jq -r .refreshToken.token <<< "$RESPONSE")
fi
echo Response : $RESPONSE
echo "-- Login admin user token end --"

# User
echo "-- Login user token start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -H "username:$USER_LOGIN_ID" -H "password:$USER_PASSWD" \
  localhost:9090 Token/LoginToken)
if [ $? != 0 ] ; then
  EXIT_CODE=1
else
  USER_ACCESS_TOKEN=$(jq -r .accessToken.token <<< "$RESPONSE")
  USER_REFRESH_TOKEN=$(jq -r .refreshToken.token <<< "$RESPONSE")
fi
echo Response : $RESPONSE
echo "-- Login user token end --"

## Refresh access token
# Admin
echo "-- Refresh admin user access token start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -d "{\"refreshToken\": \"$ADMIN_REFRESH_TOKEN\"}" \
  localhost:9090 Token/RefreshToken)
if [ $? != 0 ] ; then
  EXIT_CODE=1
fi
echo Response : $RESPONSE
echo "-- Refresh admin user access token end --"

# User
echo "-- Refresh user access token start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -d "{\"refreshToken\": \"$USER_REFRESH_TOKEN\"}" \
  localhost:9090 Token/RefreshToken)
if [ $? != 0 ] ; then
  EXIT_CODE=1
fi
echo Response : $RESPONSE
echo "-- Refresh user access token end --"

## Update user me
# Admin
echo "-- Update admin user start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -H "authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  -d "{\"phone\": \"$ADMIN_PHONE2\"}" \
  localhost:9090 UserMe/UpdateUserMe)
if [ $? != 0 ] ; then
  EXIT_CODE=1
fi
echo Response : $RESPONSE
echo "-- Update admin user End --"

# User
echo "-- Update user start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -H "authorization: Bearer $USER_ACCESS_TOKEN" \
  -d "{\"phone\": \"$USER_PHONE2\"}" \
  localhost:9090 UserMe/UpdateUserMe)
if [ $? != 0 ] ; then
  EXIT_CODE=1
fi
echo Response : $RESPONSE
echo "-- Update user End --"

## Get user me
# Admin
echo "-- Get admin user me start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -H "authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  localhost:9090 UserMe/GetUserMe)
if [ $? != 0 ] ; then
  EXIT_CODE=1
else
  ADMIN_ID=$(jq -r .id <<< "$RESPONSE")
  ADMIN_RESPONSE_PHONE=$(jq -r .phone <<< "$RESPONSE")
  if [ $ADMIN_RESPONSE_PHONE != $ADMIN_PHONE2 ]; then
    echo "!! Admin user phone info is diff !!"
    EXIT_CODE=1
  fi
fi
echo Response : $RESPONSE
echo "-- Get admin user me end --"

# User
echo "-- Get user me start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -H "authorization: Bearer $USER_ACCESS_TOKEN" \
  localhost:9090 UserMe/GetUserMe)
if [ $? != 0 ] ; then
  EXIT_CODE=1
else
  USER_ID=$(jq -r .id <<< "$RESPONSE")
  USER_RESPONSE_PHONE=$(jq -r .phone <<< "$RESPONSE")
  if [ $USER_RESPONSE_PHONE != $USER_PHONE2 ]; then
    echo "!! User user phone info is diff !!"
    EXIT_CODE=1
  fi
fi
echo Response : $RESPONSE
echo "-- Get user me end --" 
  
## Get user with admin/user token
# Get user with admin user token / Success
echo "-- Get user with admin token start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -H "authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  -d "{\"id\": \"$USER_ID\"}" \
  localhost:9090 User/GetUser)
if [ $? != 0 ] ; then
  EXIT_CODE=1
fi
echo Response : $RESPONSE
echo "-- Get user with admin token end --"

# Get admin user with user token / Fail
echo "-- Get user with user token start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -H "authorization: Bearer $USER_ACCESS_TOKEN" \
  -d "{\"id\": \"$ADMIN_ID\"}" \
  localhost:9090 User/GetUser)
if [ $? == 0 ] ; then
  EXIT_CODE=1
else 
  CODE=$(jq -r .code <<< "$RESPONSE")
  if [ $CODE != "7" ] ; then
    echo "-- Wrong code --"
    EXIT_CODE=1
  fi
fi
echo Response : $RESPONSE
echo "-- Get user with user token end --"

## Delete user me
# Admin
echo "-- Delete admin user me start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -H "authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  localhost:9090 UserMe/DeleteUserMe)
if [ $? != 0 ] ; then
  EXIT_CODE=1
fi
echo Response : $RESPONSE
echo "-- Delete admin user me end --"

# User
echo "-- Delete user me start --"
RESPONSE=$(grpcurl -plaintext -format-error \
  -H "authorization: Bearer $USER_ACCESS_TOKEN" \
  localhost:9090 UserMe/DeleteUserMe)
if [ $? != 0 ] ; then
  EXIT_CODE=1
fi
echo Response : $RESPONSE
echo "-- Delete user me end --"

# -- Exit --
exit $EXIT_CODE