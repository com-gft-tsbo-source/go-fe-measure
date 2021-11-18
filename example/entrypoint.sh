#!/usr/bin/env bash

ARGS=( )

if [[ ! -z "$MS_VERSION" ]] ;        then ARGS+=( "-version" "$MS_VERSION" ) ; fi
if [[ ! -z "$MS_PORT" ]] ;           then ARGS+=( "-port" "$MS_PORT" ) ; fi
if [[ ! -z "$MS_HOST" ]] ;           then ARGS+=( "-host" "$MS_HOST" ) ; fi
if [[ ! -z "$MS_CLIENTTIMEOUT" ]] ;  then ARGS+=( "-clienttimeout" "$MS_CLIENTTIMEOUT" ) ; fi
if [[ ! -z "$MS_CA" ]] ;             then ARGS+=( "-ca" "$MS_CA" ) ; fi
if [[ ! -z "$MS_CERT" ]] ;           then ARGS+=( "-cert" "$MS_CERT" ) ; fi
if [[ ! -z "$MS_KEY" ]] ;            then ARGS+=( "-key" "$MS_KEY" ) ; fi
if [[ ! -z "$MS_CONFIG" ]] ;         then ARGS+=( "-config" "$MS_CONFIG" ) ; fi
if [[ ! -z "$MS_DELAYREPLY" ]] ;     then ARGS+=( "-delayreply" "$MS_DELAYREPLY" ) ; fi
if [[ ! -z "$MS_MAXCONNECTIONS" ]] ; then ARGS+=( "-maxconnections" "$MS_MAXCONNECTIONS" ) ; fi
if [[ ! -z "$MS_NAME" ]] ;           then ARGS+=( "-name" "$MS_NAME" ) ; fi
if [[ ! -z "$MS_NAMESPACE" ]] ;      then ARGS+=( "-namespace" "$MS_NAMESPACE" ) ; fi
if [[ ! -z "$MS_STATICDIR" ]] ;      then ARGS+=( "-staticdir" "$MS_STATICDIR" ) ; fi
if [[ ! -z "$MS_STATICURL" ]] ;      then ARGS+=( "-staticurl" "$MS_STATICURL" ) ; fi
if [[ ! -z "$MS_TEMPLATEFILE" ]] ;      then ARGS+=( "-templatefile" "$MS_TEMPLATEFILE" ) ; fi
if [[ ! -z "$MS_TEMPLATEURL" ]] ;      then ARGS+=( "-templateurl" "$MS_TEMPLATEURL" ) ; fi

if [[ ! -z "$APP_NAME" ]] ;           then ARGS+=( "-name" "$APP_NAME" ) ; fi
if [[ ! -z "$APP_VERSION" ]] ;        then ARGS+=( "-version" "$APP_VERSION" ) ; fi


echo -n "/<NAME>"
for i in "${ARGS[@]}" "$@" ; do
    echo -n " '$i'"
done
echo

exec /<NAME> "${ARGS[@]}" "$@"