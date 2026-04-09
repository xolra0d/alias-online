package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/xolra0d/alias-online/shared/pkg/middleware"
	pbAuth "github.com/xolra0d/alias-online/shared/proto/auth"
	pbRoomManager "github.com/xolra0d/alias-online/shared/proto/room_manager"
	pbVocabManager "github.com/xolra0d/alias-online/shared/proto/vocab_manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type stubAuthClient struct {
	registerResp *pbAuth.RegisterResponse
	registerErr  error
	registerReq  *pbAuth.RegisterRequest

	loginResp *pbAuth.LoginResponse
	loginErr  error
	loginReq  *pbAuth.LoginRequest
}

func (s *stubAuthClient) Ping(context.Context, *emptypb.Empty, ...grpc.CallOption) (*pbAuth.PingResponse, error) {
	return &pbAuth.PingResponse{}, nil
}

func (s *stubAuthClient) Register(_ context.Context, req *pbAuth.RegisterRequest, _ ...grpc.CallOption) (*pbAuth.RegisterResponse, error) {
	s.registerReq = req
	return s.registerResp, s.registerErr
}

func (s *stubAuthClient) Login(_ context.Context, req *pbAuth.LoginRequest, _ ...grpc.CallOption) (*pbAuth.LoginResponse, error) {
	s.loginReq = req
	return s.loginResp, s.loginErr
}

type stubRoomManagerClient struct {
	getRoomWorkerResp *pbRoomManager.GetRoomWorkerResponse
	getRoomWorkerErr  error
	getRoomWorkerReq  *pbRoomManager.GetRoomWorkerRequest
}

func (s *stubRoomManagerClient) Ping(context.Context, *emptypb.Empty, ...grpc.CallOption) (*pbRoomManager.PingResponse, error) {
	return &pbRoomManager.PingResponse{}, nil
}

func (s *stubRoomManagerClient) GetRoomWorker(_ context.Context, req *pbRoomManager.GetRoomWorkerRequest, _ ...grpc.CallOption) (*pbRoomManager.GetRoomWorkerResponse, error) {
	s.getRoomWorkerReq = req
	return s.getRoomWorkerResp, s.getRoomWorkerErr
}

func (s *stubRoomManagerClient) PingWorker(context.Context, *pbRoomManager.PingWorkerRequest, ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *stubRoomManagerClient) RegisterRoom(context.Context, *pbRoomManager.RegisterRoomRequest, ...grpc.CallOption) (*pbRoomManager.RegisterRoomResponse, error) {
	return &pbRoomManager.RegisterRoomResponse{}, nil
}

func (s *stubRoomManagerClient) ProlongRoom(context.Context, *pbRoomManager.ProlongRoomRequest, ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *stubRoomManagerClient) ReleaseRoom(context.Context, *pbRoomManager.ReleaseRoomRequest, ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

type stubVocabManagerClient struct {
	getAvailableVocabsResp *pbVocabManager.GetAvailableVocabsResponse
	getAvailableVocabsErr  error
}

func (s *stubVocabManagerClient) Ping(context.Context, *emptypb.Empty, ...grpc.CallOption) (*pbVocabManager.PingResponse, error) {
	return &pbVocabManager.PingResponse{}, nil
}

func (s *stubVocabManagerClient) GetAvailableVocabs(context.Context, *emptypb.Empty, ...grpc.CallOption) (*pbVocabManager.GetAvailableVocabsResponse, error) {
	return s.getAvailableVocabsResp, s.getAvailableVocabsErr
}

func (s *stubVocabManagerClient) GetVocab(context.Context, *pbVocabManager.GetVocabRequest, ...grpc.CallOption) (*pbVocabManager.GetVocabResponse, error) {
	return &pbVocabManager.GetVocabResponse{}, nil
}

func newTestHandles(auth pbAuth.AuthServiceClient, vocab pbVocabManager.VocabManagerServiceClient, room pbRoomManager.RoomManagerServiceClient) *Handles {
	return NewHTTPHandles(
		auth,
		vocab,
		room,
		testLogger(),
		10*time.Second,
		10*time.Second,
		time.Hour,
		"/",
		false,
		true,
		"localhost",
	)
}

func decodeJSONBody(t *testing.T, rr *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var got map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode body as json: %v; body=%q", err, rr.Body.String())
	}
	return got
}

func TestIsValidRoomId(t *testing.T) {
	cases := []struct {
		roomId string
		want   bool
	}{
		{roomId: "ABCDEF", want: true},
		{roomId: "AZ2345", want: true},
		{roomId: "ABCDE", want: false},
		{roomId: "ABCDEFG", want: false},
		{roomId: "abcDEF", want: false},
		{roomId: "ABCD18", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.roomId, func(t *testing.T) {
			if got := isValidRoomId(tc.roomId); got != tc.want {
				t.Fatalf("isValidRoomId(%q)=%v, want %v", tc.roomId, got, tc.want)
			}
		})
	}
}

func TestAvailableVocabsSuccess(t *testing.T) {
	vocabClient := &stubVocabManagerClient{
		getAvailableVocabsResp: &pbVocabManager.GetAvailableVocabsResponse{Names: []string{"animals", "cities"}},
	}
	h := newTestHandles(&stubAuthClient{}, vocabClient, &stubRoomManagerClient{})

	req := httptest.NewRequest(http.MethodGet, "/api/available-vocabs", nil)
	rr := httptest.NewRecorder()
	h.AvailableVocabs(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	got := decodeJSONBody(t, rr)
	vocabs, ok := got["vocabs"].([]any)
	if !ok {
		t.Fatalf("vocabs is not a JSON array: %#v", got["vocabs"])
	}
	if len(vocabs) != 2 || vocabs[0] != "animals" || vocabs[1] != "cities" {
		t.Fatalf("unexpected vocabs payload: %#v", vocabs)
	}
}

func TestRegisterSetsCookieAndReturnsOK(t *testing.T) {
	authClient := &stubAuthClient{
		registerResp: &pbAuth.RegisterResponse{
			Token: "jwt-token",
			Exp:   time.Now().Add(time.Hour).Unix(),
		},
	}
	h := newTestHandles(authClient, &stubVocabManagerClient{}, &stubRoomManagerClient{})

	body := bytes.NewBufferString(`{"name":"username1","login":"loginname","password":"password1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/register", body)
	rr := httptest.NewRecorder()
	h.Register(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if authClient.registerReq == nil {
		t.Fatal("register request was not sent to auth client")
	}

	got := decodeJSONBody(t, rr)
	if ok, isBool := got["ok"].(bool); !isBool || !ok {
		t.Fatalf("unexpected response payload: %#v", got)
	}

	cookies := rr.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected exactly one cookie, got %d", len(cookies))
	}
	c := cookies[0]
	if c.Name != middleware.LoginCookieName {
		t.Fatalf("cookie name=%q, want %q", c.Name, middleware.LoginCookieName)
	}
	if c.Value != "jwt-token" {
		t.Fatalf("cookie value=%q, want %q", c.Value, "jwt-token")
	}
	if c.Path != "/" {
		t.Fatalf("cookie path=%q, want %q", c.Path, "/")
	}
	if !c.HttpOnly {
		t.Fatal("cookie should be HttpOnly")
	}
}

func TestRegisterAlreadyExistsMapsToConflict(t *testing.T) {
	authClient := &stubAuthClient{
		registerErr: status.Error(codes.AlreadyExists, "already exists"),
	}
	h := newTestHandles(authClient, &stubVocabManagerClient{}, &stubRoomManagerClient{})

	body := bytes.NewBufferString(`{"name":"username1","login":"loginname","password":"password1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/register", body)
	rr := httptest.NewRecorder()
	h.Register(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("status=%d, want %d; body=%s", rr.Code, http.StatusConflict, rr.Body.String())
	}
	got := decodeJSONBody(t, rr)
	if got["err"] != "user already exists" {
		t.Fatalf("unexpected error payload: %#v", got)
	}
}

func TestPlayInvalidRoomIdReturnsBadRequest(t *testing.T) {
	roomClient := &stubRoomManagerClient{}
	h := newTestHandles(&stubAuthClient{}, &stubVocabManagerClient{}, roomClient)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/protected/play/{roomId}", h.Play)

	req := httptest.NewRequest(http.MethodGet, "/api/protected/play/bad-room", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want %d; body=%s", rr.Code, http.StatusBadRequest, rr.Body.String())
	}
	if roomClient.getRoomWorkerReq != nil {
		t.Fatal("room manager client should not be called for invalid room id")
	}
	got := decodeJSONBody(t, rr)
	if got["err"] != "invalid room id" {
		t.Fatalf("unexpected error payload: %#v", got)
	}
}

func TestPlaySuccessReturnsWorker(t *testing.T) {
	roomClient := &stubRoomManagerClient{
		getRoomWorkerResp: &pbRoomManager.GetRoomWorkerResponse{Worker: "ws://room-worker:8082"},
	}
	h := newTestHandles(&stubAuthClient{}, &stubVocabManagerClient{}, roomClient)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/protected/play/{roomId}", h.Play)

	req := httptest.NewRequest(http.MethodGet, "/api/protected/play/ABC234", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if roomClient.getRoomWorkerReq == nil || roomClient.getRoomWorkerReq.RoomId != "ABC234" {
		t.Fatalf("unexpected room request: %#v", roomClient.getRoomWorkerReq)
	}
	got := decodeJSONBody(t, rr)
	if got["worker"] != "ws://room-worker:8082" {
		t.Fatalf("unexpected worker payload: %#v", got)
	}
}

