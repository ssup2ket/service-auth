#!/bin/bash
# -- Variables --
EXIT_CODE=0
RANDOM_VALUE=$RANDOM

ADMIN_ID=""
ADMIN_TOKEN=""
ADMIN_LOGIN_ID="admin$RANDOM_VALUE"
ADMIN_PASSWD="admin$RANDOM_VALUE"
ADMIN_PHONE="000-0000-0000"
ADMIN_PHONE2="111-1111-1111"
ADMIN_EMAIL="admin$RANDOM_VALUE@test.com"

USER_ID=""
USER_TOKEN=""
USER_LOGIN_ID="user$RANDOM_VALUE"
USER_PASSWD="user$RANDOM_VALUE"
USER_PHONE="222-2222-2222"
USER_PHONE2="333-3333-3333"
USER_EMAIL="user$RANDOM_VALUE@test.com"

# -- Test cases --
# Create users
## Admin
echo "-- Create admin user start --"
RESPONSE=$(curl --no-progress-meter --write-out '%{http_code}' \
  -X POST \
  -H "Content-Type: application/json" \
  -d "{\"loginId\": \"$ADMIN_LOGIN_ID\", \"password\": \"$ADMIN_PASSWD\", \"phone\": \"$ADMIN_PHONE\", \"email\": \"$ADMIN_EMAIL\", \"role\" : \"admin\"}" \
  http://localhost/v1/users)
RESPONSE_HTTP_CODE=$(tail -n1 <<< "$RESPONSE")
RESPONSE_BODY=$(sed '$ d' <<< "$RESPONSE")
echo Response HTTP Code : $RESPONSE_HTTP_CODE
echo Response Body : $RESPONSE_BODY
echo "-- Create admin user end --"

## User
echo "-- Create user start --"
RESPONSE=$(curl --no-progress-meter --write-out '%{http_code}' \
  -X POST \
  -H "Content-Type: application/json" \
  -d "{\"loginId\": \"$USER_LOGIN_ID\", \"password\": \"$USER_PASSWD\", \"phone\": \"$USER_PHONE\", \"email\": \"$USER_EMAIL\", \"role\" : \"user\"}" \
  http://localhost/v1/users)
RESPONSE_HTTP_CODE=$(tail -n1 <<< "$RESPONSE")
RESPONSE_BODY=$(sed '$ d' <<< "$RESPONSE")
echo Response HTTP Code : $RESPONSE_HTTP_CODE
echo Response Body : $RESPONSE_BODY
echo "-- Create user end --"

# Create tokens
## Admin
echo "-- Create admin user token start --"
RESPONSE=$(curl --no-progress-meter --write-out '%{http_code}' \
  -X POST \
  -H "Content-Type: application/json" \
  -d "{\"loginId\": \"$ADMIN_LOGIN_ID\", \"password\": \"$ADMIN_PASSWD\"}" \
  http://localhost/v1/tokens)
RESPONSE_HTTP_CODE=$(tail -n1 <<< "$RESPONSE")
RESPONSE_BODY=$(sed '$ d' <<< "$RESPONSE")
echo Response HTTP Code : $RESPONSE_HTTP_CODE
echo Response Body : $RESPONSE_BODY
if [ $RESPONSE_HTTP_CODE != "200" ]; then
  EXIT_CODE=1
else
  ADMIN_TOKEN=$(jq -r .token <<< "$RESPONSE_BODY")
fi
echo "-- Create admin user token End --"

## User
echo "-- Create user token start --"
RESPONSE=$(curl --no-progress-meter --write-out '%{http_code}' \
  -X POST \
  -H "Content-Type: application/json" \
  -d "{\"loginId\": \"$USER_LOGIN_ID\", \"password\": \"$USER_PASSWD\"}" \
  http://localhost/v1/tokens)
RESPONSE_HTTP_CODE=$(tail -n1 <<< "$RESPONSE")
RESPONSE_BODY=$(sed '$ d' <<< "$RESPONSE")
echo Response HTTP Code : $RESPONSE_HTTP_CODE
echo Response Body : $RESPONSE_BODY
if [ $RESPONSE_HTTP_CODE != "200" ]; then
  EXIT_CODE=1
else
  USER_TOKEN=$(jq -r .token <<< "$RESPONSE_BODY")
fi
echo "-- Create user token end --"

# Update user me
## Admin
echo "-- Update admin user start --"
RESPONSE=$(curl --no-progress-meter --write-out '%{http_code}' \
  -X PUT \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "{\"phone\": \"$ADMIN_PHONE2\"}" \
  http://localhost/v1/users/me)
RESPONSE_HTTP_CODE=$(tail -n1 <<< "$RESPONSE")
echo Response HTTP Code : $RESPONSE_HTTP_CODE
if [ $RESPONSE_HTTP_CODE != "200" ]; then
  EXIT_CODE=1
fi
echo "-- Update admin user End --"

## User
echo "-- Update user start --"
RESPONSE=$(curl --no-progress-meter --write-out '%{http_code}' \
  -X PUT \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d "{\"phone\": \"$USER_PHONE2\"}" \
  http://localhost/v1/users/me)
RESPONSE_HTTP_CODE=$(tail -n1 <<< "$RESPONSE")
echo Response HTTP Code : $RESPONSE_HTTP_CODE
if [ $RESPONSE_HTTP_CODE != "200" ]; then
  EXIT_CODE=1
fi
echo "-- Update user End --"

# Get user me
## Admin
echo "-- Get admin user me start --"
RESPONSE=$(curl --no-progress-meter --write-out '%{http_code}' \
  -X GET \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "{\"loginId\": \"$ADMIN_LOGIN_ID\", \"password\": \"$ADMIN_PASSWD\"}" \
  http://localhost/v1/users/me)
RESPONSE_HTTP_CODE=$(tail -n1 <<< "$RESPONSE")
RESPONSE_BODY=$(sed '$ d' <<< "$RESPONSE")
echo Response HTTP Code : $RESPONSE_HTTP_CODE
echo Response Body : $RESPONSE_BODY
if [ $RESPONSE_HTTP_CODE != "200" ]; then
  EXIT_CODE=1
else
  ADMIN_ID=$(jq -r .id <<< "$RESPONSE_BODY")
  ADMIN_RESPONSE_PHONE=$(jq -r .phone <<< "$RESPONSE_BODY")
  if [$ADMIN_RESPONSE_PHONE != $ADMIN_PHONE2]; then
    echo "!! Admin user phone info is diff !!"
    EXIT_CODE=1
  fi
fi
echo "-- Get admin user me end --"

## User
echo "-- Get user me start --"
RESPONSE=$(curl --no-progress-meter --write-out '%{http_code}' \
  -X GET \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d "{\"loginId\": \"$USER_LOGIN_ID\", \"password\": \"$USER_PASSWD\"}" \
  http://localhost/v1/users/me)
RESPONSE_HTTP_CODE=$(tail -n1 <<< "$RESPONSE")
RESPONSE_BODY=$(sed '$ d' <<< "$RESPONSE")
echo Response HTTP Code : $RESPONSE_HTTP_CODE
echo Response Body : $RESPONSE_BODY
if [ $RESPONSE_HTTP_CODE != "200" ]; then
  EXIT_CODE=1
else
  USER_ID=$(jq -r .id <<< "$RESPONSE_BODY")
  USER_RESPONSE_PHONE=$(jq -r .phone <<< "$RESPONSE_BODY")
  if [$USER_RESPONSE_PHONE != $USER_PHONE2]; then
    echo "!! User phone info is diff !!"
    EXIT_CODE=1
  fi
fi
echo "-- Get user me end --"

# Delete user me
## Admin
echo "-- Delete admin user me start --"
RESPONSE=$(curl --no-progress-meter --write-out '%{http_code}' \
  -X DELETE \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost/v1/users/me)
RESPONSE_HTTP_CODE=$(tail -n1 <<< "$RESPONSE")
echo Response HTTP Code : $RESPONSE_HTTP_CODE
if [ $RESPONSE_HTTP_CODE != "200" ]; then
  EXIT_CODE=1
fi
echo "-- Delete admin user me end --"

## User
echo "-- Delete user me start --"
RESPONSE=$(curl --no-progress-meter --write-out '%{http_code}' \
  -X DELETE \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  http://localhost/v1/users/me)
RESPONSE_HTTP_CODE=$(tail -n1 <<< "$RESPONSE")
echo Response HTTP Code : $RESPONSE_HTTP_CODE
if [ $RESPONSE_HTTP_CODE != "200" ]; then
  EXIT_CODE=1
fi
echo "-- Delete user me end --"

# -- Exit --
exit $EXIT_CODE