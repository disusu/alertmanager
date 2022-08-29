{-
   Alertmanager API
   API of the Prometheus Alertmanager (https://github.com/prometheus/alertmanager)

   OpenAPI spec version: 0.0.1

   NOTE: This file is auto generated by the openapi-generator.
   https://github.com/openapitools/openapi-generator.git
   Do not edit this file manually.
-}


module Data.PostableSlsAlerts exposing (PostableSlsAlerts, decoder, encoder)

import DateTime exposing (DateTime)
import Dict exposing (Dict)
import Json.Decode as Decode exposing (Decoder)
import Json.Decode.Pipeline exposing (optional, required)
import Json.Encode as Encode


type alias PostableSlsAlerts =
    { startsAt : Maybe DateTime
    , endsAt : Maybe DateTime
    , results : Maybe (List (Dict String String))
    , annotations : Maybe (Dict String String)
    , generatorURL : Maybe String
    }


decoder : Decoder PostableSlsAlerts
decoder =
    Decode.succeed PostableSlsAlerts
        |> optional "startsAt" (Decode.nullable DateTime.decoder) Nothing
        |> optional "endsAt" (Decode.nullable DateTime.decoder) Nothing
        |> optional "results" (Decode.nullable (Decode.list (Decode.dict Decode.string))) Nothing
        |> optional "annotations" (Decode.nullable (Decode.dict Decode.string)) Nothing
        |> optional "generatorURL" (Decode.nullable Decode.string) Nothing


encoder : PostableSlsAlerts -> Encode.Value
encoder model =
    Encode.object
        [ ( "startsAt", Maybe.withDefault Encode.null (Maybe.map DateTime.encoder model.startsAt) )
        , ( "endsAt", Maybe.withDefault Encode.null (Maybe.map DateTime.encoder model.endsAt) )
        , ( "results", Maybe.withDefault Encode.null (Maybe.map (Encode.list (Encode.dict identity Encode.string)) model.results) )
        , ( "annotations", Maybe.withDefault Encode.null (Maybe.map (Encode.dict identity Encode.string) model.annotations) )
        , ( "generatorURL", Maybe.withDefault Encode.null (Maybe.map Encode.string model.generatorURL) )
        ]
