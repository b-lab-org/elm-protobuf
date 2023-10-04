module Recursive exposing (..)

-- DO NOT EDIT
-- AUTOGENERATED BY THE ELM PROTOCOL BUFFER COMPILER
-- https://github.com/tiziano88/elm-protobuf
-- source file: recursive.proto

import Protobuf exposing (..)

import Json.Decode as JD
import Json.Encode as JE


uselessDeclarationToPreventErrorDueToEmptyOutputFile = 42


type alias Rec =
    { int32Field : Int -- 1
    , stringField : String -- 4
    , r : Rec_R
    }


recDecoder : JD.Decoder Rec
recDecoder =
    JD.lazy <| \_ -> decode Rec
        |> required "int32Field" intDecoder 0
        |> required "stringField" JD.string ""
        |> field rec_RDecoder


recEncoder : Rec -> JE.Value
recEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "int32Field" JE.int 0 v.int32Field)
        , (requiredFieldEncoder "stringField" JE.string "" v.stringField)
        , (rec_REncoder v.r)
        ]


type Rec_R
    = Rec_RUnspecified
    | RecField Rec


rec_RDecoder : JD.Decoder Rec_R
rec_RDecoder =
    JD.lazy <| \_ -> JD.oneOf
        [ JD.map RecField (JD.field "recField" recDecoder)
        , JD.succeed Rec_RUnspecified
        ]


rec_REncoder : Rec_R -> Maybe ( String, JE.Value )
rec_REncoder v =
    case v of
        Rec_RUnspecified ->
            Nothing

        RecField x ->
            Just ( "recField", recEncoder x )
