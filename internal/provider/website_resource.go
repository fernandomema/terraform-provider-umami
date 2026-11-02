package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fernandomema/terraform-provider-umami/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &websiteResource{}
	_ resource.ResourceWithConfigure   = &websiteResource{}
	_ resource.ResourceWithImportState = &websiteResource{}
)

func NewWebsiteResource() resource.Resource { return &websiteResource{} }

type websiteResource struct{ data *providerData }

type websiteModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Domain         types.String `tfsdk:"domain"`
	ShareID        types.String `tfsdk:"share_id"`
	ScriptURL      types.String `tfsdk:"script_url"`
	TrackingScript types.String `tfsdk:"tracking_script"`
}

func (r *websiteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_website"
}

func (r *websiteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Un sitio (website) de Umami. Su id es el data-website-id del tracker.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "ID del sitio = data-website-id.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name":     schema.StringAttribute{Required: true, Description: "Nombre visible."},
			"domain":   schema.StringAttribute{Required: true, Description: "Dominio del sitio (sin protocolo)."},
			"share_id": schema.StringAttribute{Optional: true, Description: "ID para compartir el dashboard publicamente (/share/<id>)."},
			"script_url": schema.StringAttribute{
				Computed:      true,
				Description:   "URL del script.js de la instancia.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tracking_script": schema.StringAttribute{
				Computed:      true,
				Description:   "Etiqueta <script> lista para pegar en el sitio.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *websiteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d, ok := req.ProviderData.(*providerData)
	if !ok {
		resp.Diagnostics.AddError("Provider data inesperado", fmt.Sprintf("%T", req.ProviderData))
		return
	}
	r.data = d
}

func (r *websiteResource) input(m websiteModel) client.WebsiteInput {
	in := client.WebsiteInput{Name: m.Name.ValueString(), Domain: m.Domain.ValueString()}
	if !m.ShareID.IsNull() && m.ShareID.ValueString() != "" {
		s := m.ShareID.ValueString()
		in.ShareID = &s
	}
	return in
}

func (r *websiteResource) fill(m *websiteModel, w *client.Website) {
	m.ID = types.StringValue(w.ID)
	m.Name = types.StringValue(w.Name)
	m.Domain = types.StringValue(w.Domain)
	if w.ShareID != nil && *w.ShareID != "" {
		m.ShareID = types.StringValue(*w.ShareID)
	} else {
		m.ShareID = types.StringNull()
	}
	script := strings.TrimRight(r.data.Endpoint, "/") + "/script.js"
	m.ScriptURL = types.StringValue(script)
	m.TrackingScript = types.StringValue(fmt.Sprintf(`<script defer src="%s" data-website-id="%s"></script>`, script, w.ID))
}

func (r *websiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan websiteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	w, err := r.data.Client.CreateWebsite(ctx, r.input(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error creando el website", err.Error())
		return
	}
	r.fill(&plan, w)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *websiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state websiteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	w, err := r.data.Client.GetWebsite(ctx, state.ID.ValueString())
	var nf *client.NotFoundError
	if errors.As(err, &nf) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error leyendo el website", err.Error())
		return
	}
	r.fill(&state, w)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *websiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state websiteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	w, err := r.data.Client.UpdateWebsite(ctx, state.ID.ValueString(), r.input(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error actualizando el website", err.Error())
		return
	}
	r.fill(&plan, w)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *websiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state websiteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.data.Client.DeleteWebsite(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error borrando el website", err.Error())
	}
}

func (r *websiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
